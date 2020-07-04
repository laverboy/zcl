package zcl

import (
	"fmt"
	"github.com/shimmeringbee/zigbee"
	"reflect"
)

type CommandRegistry struct {
	globalIdentifierToInterface map[CommandIdentifier]interface{}
	globalInterfaceToIdentifier map[reflect.Type]CommandIdentifier

	localIdentifierToInterface map[zigbee.ClusterID]map[zigbee.ManufacturerCode]map[CommandIdentifier]interface{}
	localInterfaceToIdentifier map[zigbee.ClusterID]map[zigbee.ManufacturerCode]map[reflect.Type]CommandIdentifier
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		globalIdentifierToInterface: make(map[CommandIdentifier]interface{}),
		globalInterfaceToIdentifier: make(map[reflect.Type]CommandIdentifier),
		localIdentifierToInterface:  make(map[zigbee.ClusterID]map[zigbee.ManufacturerCode]map[CommandIdentifier]interface{}),
		localInterfaceToIdentifier:  make(map[zigbee.ClusterID]map[zigbee.ManufacturerCode]map[reflect.Type]CommandIdentifier),
	}
}

func (cr *CommandRegistry) RegisterGlobal(identifier CommandIdentifier, command interface{}) {
	cr.globalIdentifierToInterface[identifier] = command
	cr.globalInterfaceToIdentifier[reflect.TypeOf(command)] = identifier
}

func (cr *CommandRegistry) GetGlobalCommand(identifier CommandIdentifier) (interface{}, error) {
	sampleObj, found := cr.globalIdentifierToInterface[identifier]

	if found {
		reflectedType := reflect.TypeOf(sampleObj).Elem()
		return reflect.New(reflectedType).Interface(), nil
	} else {
		return 0, fmt.Errorf("could not find global command for identifier: %d", identifier)
	}
}

func (cr *CommandRegistry) GetGlobalCommandIdentifier(command interface{}) (CommandIdentifier, error) {
	reflectedType := reflect.TypeOf(command)
	identifier, found := cr.globalInterfaceToIdentifier[reflectedType]

	if found {
		return identifier, nil
	} else {
		return 0, fmt.Errorf("could not find global command for type: %s", reflectedType.Name())
	}
}

func (cr *CommandRegistry) RegisterLocal(clusterID zigbee.ClusterID, manufacturer zigbee.ManufacturerCode, identifier CommandIdentifier, command interface{}) {
	clusterId2IntResult, clusterId2IntFound := cr.localIdentifierToInterface[clusterID]
	clusterInt2IdResult, clusterInt2IdFound := cr.localInterfaceToIdentifier[clusterID]

	if !clusterId2IntFound {
		clusterId2IntResult = make(map[zigbee.ManufacturerCode]map[CommandIdentifier]interface{})
		cr.localIdentifierToInterface[clusterID] = clusterId2IntResult
	}

	if !clusterInt2IdFound {
		clusterInt2IdResult = make(map[zigbee.ManufacturerCode]map[reflect.Type]CommandIdentifier)
		cr.localInterfaceToIdentifier[clusterID] = clusterInt2IdResult
	}

	manufacturerId2IntResult, manufacturerId2IntFound := clusterId2IntResult[manufacturer]
	manufacturerInt2IdResult, manufacturerInt2IdFound := clusterInt2IdResult[manufacturer]

	if !manufacturerId2IntFound {
		manufacturerId2IntResult = make(map[CommandIdentifier]interface{})
		clusterId2IntResult[manufacturer] = manufacturerId2IntResult
	}

	if !manufacturerInt2IdFound {
		manufacturerInt2IdResult = make(map[reflect.Type]CommandIdentifier)
		clusterInt2IdResult[manufacturer] = manufacturerInt2IdResult
	}

	manufacturerId2IntResult[identifier] = command
	manufacturerInt2IdResult[reflect.TypeOf(command)] = identifier
}

func (cr *CommandRegistry) GetLocalCommand(clusterID zigbee.ClusterID, manufacturer zigbee.ManufacturerCode, identifier CommandIdentifier) (interface{}, error) {
	clusterResult, clusterFound := cr.localIdentifierToInterface[clusterID]

	if !clusterFound {
		return nil, fmt.Errorf("could not find local command for: cluster %d manufacturer %d identifier %d", clusterID, manufacturer, identifier)
	}

	manufacturerResult, manufacturerFound := clusterResult[manufacturer]

	if !manufacturerFound {
		return nil, fmt.Errorf("could not find local command for: cluster %d manufacturer %d identifier %d", clusterID, manufacturer, identifier)
	}

	interfaceResult, interfaceFound := manufacturerResult[identifier]

	if !interfaceFound {
		return nil, fmt.Errorf("could not find local command for: cluster %d manufacturer %d identifier %d", clusterID, manufacturer, identifier)
	}

	reflectedType := reflect.TypeOf(interfaceResult).Elem()
	return reflect.New(reflectedType).Interface(), nil
}

func (cr *CommandRegistry) GetLocalCommandIdentifier(clusterID zigbee.ClusterID, manufacturer zigbee.ManufacturerCode, command interface{}) (CommandIdentifier, error) {
	reflectedType := reflect.TypeOf(command)
	clusterResult, clusterFound := cr.localInterfaceToIdentifier[clusterID]

	if !clusterFound {
		return 0, fmt.Errorf("could not find local command for: cluster %d manufacturer %d type %s", clusterID, manufacturer, reflectedType.Name())
	}

	manufacturerResult, manufacturerFound := clusterResult[manufacturer]

	if !manufacturerFound {
		return 0, fmt.Errorf("could not find local command for: cluster %d manufacturer %d type %s", clusterID, manufacturer, reflectedType.Name())
	}

	identifierResult, identifierFound := manufacturerResult[reflectedType]

	if !identifierFound {
		return 0, fmt.Errorf("could not find local command for: cluster %d manufacturer %d type %s", clusterID, manufacturer, reflectedType.Name())
	}

	return identifierResult, nil
}
