package cm

import (
	"bytes"
	"html/template"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/utils"
	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/volume"
	corev1 "k8s.io/api/core/v1"
)

var (
	ConfigVolume = SharedVolume{
		Name:                   "zookeeper-internal-config-local",
		InitContainerMountPath: "/mnt/zookeeper/zookeeper-config-local",
		ContainerMountPath:     volume.ConfigVolumeMountPath,
	}

	PluginVolumes = SharedVolumeArray{
		Array: []SharedVolume{
			ConfigVolume,
		},
	}

	// linkedFiles describe how various secrets are mapped into the pod's filesystem.
	linkedFiles = LinkedFilesArray{
		Array: []LinkedFile{
			{
				Source: utils.Joins(ConfigVolumeMountPath, "/", ConfigFileName),
				Target: utils.Joins(ConfigVolume.ContainerMountPath, "/", ConfigFileName),
			},
			{
				Source: utils.Joins(DynamicConfigFileVolumeMountPath, "/", DynamicConfigFile),
				Target: utils.Joins(ConfigVolume.ContainerMountPath, "/", DynamicConfigFile),
			},
		},
	}
)

func (v SharedVolumeArray) Volumes() []corev1.Volume {
	volumes := make([]corev1.Volume, len(v.Array))
	for i, v := range v.Array {
		volumes[i] = corev1.Volume{
			Name: v.Name,
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
	}
	return volumes
}

// SharedVolumes represents a list of SharedVolume
type SharedVolumeArray struct {
	Array []SharedVolume
}

// SharedVolume between the init container and the ZK container.
type SharedVolume struct {
	Name                   string // Volume name
	InitContainerMountPath string // Mount path in the init container
	ContainerMountPath     string // Mount path in the zookeeper container
}

// LinkedFilesArray contains all files to be linked in the init container.
type LinkedFilesArray struct {
	Array []LinkedFile
}

// LinkedFile describes a symbolic link with source and target.
type LinkedFile struct {
	Source string
	Target string
}

// TemplateParams are the parameters manipulated in the scriptTemplate
type TemplateParams struct {
	// SharedVolumes are directories to persist in shared volumes
	PluginVolumes SharedVolumeArray
	// LinkedFiles are files to link individually
	LinkedFiles LinkedFilesArray
	// ChownToElasticsearch are paths that need to be chowned to the Elasticsearch user/group.
	ChownToZookeeper []string
}

// RenderScriptTemplate renders scriptTemplate using the given TemplateParams
func RenderScriptTemplate(params TemplateParams) (string, error) {
	tplBuffer := bytes.Buffer{}
	//fmt.Println(params)
	if err := scriptTemplate.Execute(&tplBuffer, params); err != nil {
		return "", err
	}
	return tplBuffer.String(), nil
}

const (
	PrepareFsScriptConfigKey = "prepare-fs.sh"
)

// scriptTemplate is the main script to be run
// in the prepare-fs init container before ES starts
var scriptTemplate = template.Must(template.New("").Parse(
	`#!/usr/bin/env bash

	set -eu

	# compute time in seconds since the given start time
	function duration() {
		local start=$1
		end=$(date +%s)
		echo $((end-start))
	}

	######################
	#        START       #
	######################

	script_start=$(date +%s)

	echo "Starting init script"

	######################
	#  Config linking    #
	######################

	# Link individual files from their mount location into the config dir
	# to a volume, to be used by the ZK container
	ln_start=$(date +%s)
	{{range .LinkedFiles.Array}}
		echo "Linking {{.Source}} to {{.Target}}"
		ln -sf {{.Source}} {{.Target}}
	{{end}}
	ls -l /conf
	echo "File linking duration: $(duration $ln_start) sec."

	######################
	#  Volumes chown     #
	######################

	# chown the data and logs volume to the elasticsearch user
	# only done when running as root, other cases should be handled
	# with a proper security context
	chown_start=$(date +%s)
	if [[ $EUID -eq 0 ]]; then
		{{range .ChownToZookeeper}}
			echo "chowning {{.}} to zookeeper:zookeeper"
			chown -v zookeeper:zookeeper {{.}}
		{{end}}
	fi
	echo "chown duration: $(duration $chown_start) sec."

	######################
	#         End        #
	######################

	echo "Init script successful"
	echo "Script duration: $(duration $script_start) sec."
`))

func RenderPrepareFsScript() (string, error) {
	return RenderScriptTemplate(TemplateParams{
		PluginVolumes: PluginVolumes,
		LinkedFiles:   linkedFiles,
		ChownToZookeeper: []string{
			DataMountPath,
			LogsMountPath,
		},
	})
}
