package repositories

import (
	"errors"
	"fmt"
	"net/http"
	"orchestration_service/models"
	"slices"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceRepository struct {
	devices []*models.Device
	log     *logr.Logger
}

func (r *DeviceRepository) AddDevice(m *models.DeviceDescription) (*models.DeviceDescription, *models.ResponseError) {
	r.log.Info(fmt.Sprintf("Adding device: id='%s' name='%s'", m.ID, m.Metadata.Name))
	idx := slices.IndexFunc(r.devices, func(d *models.Device) bool { return d.Description.ID == m.ID })
	if idx >= 0 {
		r.log.Info("Device already exists")
		return &r.devices[idx].Description, nil
	}

	r.devices = append(r.devices, &models.Device{Description: *m})
	r.log.Info(fmt.Sprint("Added device: ", &r.devices))
	return m, nil
}

func (r *DeviceRepository) GetDevice(id uuid.UUID) (*models.DeviceDescription, *models.ResponseError) {
	r.log.Info(fmt.Sprint("Deleting devices ", id))
	idx := slices.IndexFunc(r.devices, func(d *models.Device) bool { return d.Description.ID == id })
	if idx < 0 {
		msg := fmt.Sprint("Unable to locate device ", id)
		r.log.Error(errors.New("device not found"), msg)

		return nil, &models.ResponseError{
			Error:  msg,
			Status: http.StatusNotFound,
		}
	}
	r.log.Info(fmt.Sprint("Found device: ", &r.devices[idx].Description))
	return &r.devices[idx].Description, nil
}

func (r *DeviceRepository) GetDevices() ([]models.DeviceDescription, *models.ResponseError) {
	r.log.Info("Getting devices")
	d := make([]models.DeviceDescription, len(r.devices))
	for i, v := range r.devices {
		d[i] = v.Description
	}

	r.log.Info(fmt.Sprint("Found devices: ", d))
	return d, nil
}

func (r *DeviceRepository) AddApp(deviceId uuid.UUID, appId uuid.UUID) *models.ResponseError {
	r.log.Info(fmt.Sprintf("Adding app %s to device %s", appId, deviceId))
	idx := slices.IndexFunc(r.devices, func(d *models.Device) bool { return d.Description.ID == deviceId })
	if idx < 0 {
		msg := fmt.Sprint("Unable to locate device ", deviceId)
		r.log.Error(errors.New("device not found"), msg)

		return &models.ResponseError{
			Error:  msg,
			Status: http.StatusNotFound,
		}
	}

	r.devices[idx].Apps = append(r.devices[idx].Apps, appId)
	return nil
}

func NewDeviceRepository(l *logr.Logger) *DeviceRepository {
	return &DeviceRepository{
		log: l,
	}
}
