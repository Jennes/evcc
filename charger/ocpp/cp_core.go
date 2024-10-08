package ocpp

import (
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidConnector   = errors.New("invalid connector")
	ErrInvalidTransaction = errors.New("invalid transaction")
)

func (cp *CP) Authorize(request *core.AuthorizeRequest) (*core.AuthorizeConfirmation, error) {
	res := &core.AuthorizeConfirmation{
		IdTagInfo: &types.IdTagInfo{
			Status: types.AuthorizationStatusAccepted,
		},
	}

	return res, nil
}

func (cp *CP) BootNotification(request *core.BootNotificationRequest) (*core.BootNotificationConfirmation, error) {
	res := &core.BootNotificationConfirmation{
		CurrentTime: types.Now(),
		Interval:    60,
		Status:      core.RegistrationStatusAccepted,
	}

	cp.onceBoot.Do(func() {
		cp.bootNotificationRequestC <- request
	})

	return res, nil
}

func (cp *CP) DiagnosticStatusNotification(request *firmware.DiagnosticsStatusNotificationRequest) (*firmware.DiagnosticsStatusNotificationConfirmation, error) {
	return new(firmware.DiagnosticsStatusNotificationConfirmation), nil
}

func (cp *CP) FirmwareStatusNotification(request *firmware.FirmwareStatusNotificationRequest) (*firmware.FirmwareStatusNotificationConfirmation, error) {
	return new(firmware.FirmwareStatusNotificationConfirmation), nil
}

func (cp *CP) StatusNotification(request *core.StatusNotificationRequest) (*core.StatusNotificationConfirmation, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if conn := cp.connectorByID(request.ConnectorId); conn != nil {
		return conn.StatusNotification(request)
	}

	return new(core.StatusNotificationConfirmation), nil
}

func (cp *CP) DataTransfer(request *core.DataTransferRequest) (*core.DataTransferConfirmation, error) {
	res := &core.DataTransferConfirmation{
		Status: core.DataTransferStatusAccepted,
	}

	return res, nil
}

func (cp *CP) Heartbeat(request *core.HeartbeatRequest) (*core.HeartbeatConfirmation, error) {
	res := &core.HeartbeatConfirmation{
		CurrentTime: types.Now(),
	}

	return res, nil
}

func (cp *CP) MeterValues(request *core.MeterValuesRequest) (*core.MeterValuesConfirmation, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	// signal received
	select {
	case cp.meterC <- struct{}{}:
	default:
	}

	if conn := cp.connectorByID(request.ConnectorId); conn != nil {
		conn.MeterValues(request)
	}

	return new(core.MeterValuesConfirmation), nil
}

func (cp *CP) StartTransaction(request *core.StartTransactionRequest) (*core.StartTransactionConfirmation, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if conn := cp.connectorByID(request.ConnectorId); conn != nil {
		return conn.StartTransaction(request)
	}

	return new(core.StartTransactionConfirmation), nil
}

func (cp *CP) StopTransaction(request *core.StopTransactionRequest) (*core.StopTransactionConfirmation, error) {
	if request == nil {
		return nil, ErrInvalidRequest
	}

	if conn := cp.connectorByTransactionID(request.TransactionId); conn != nil {
		return conn.StopTransaction(request)
	}

	res := &core.StopTransactionConfirmation{
		IdTagInfo: &types.IdTagInfo{
			Status: types.AuthorizationStatusAccepted, // accept old pending stop message during startup
		},
	}

	return res, nil
}
