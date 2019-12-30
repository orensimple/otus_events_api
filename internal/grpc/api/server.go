package api

import (
	"context"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/orensimple/otus_events_api/internal/domain/errors"
	"github.com/orensimple/otus_events_api/internal/domain/services"
)

type CalendarServer struct {
	EventService *services.EventService
}

func (cs *CalendarServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {

	var startTime *time.Time
	if req.GetStartTime() != nil {
		st, err := ptypes.Timestamp(req.GetStartTime())
		if err != nil {
			return nil, err
		}
		startTime = &st
	}

	var endTime *time.Time
	if req.GetEndTime() != nil {
		et, err := ptypes.Timestamp(req.GetEndTime())
		if err != nil {
			return nil, err
		}
		endTime = &et
	}

	event, err := cs.EventService.CreateEvent(ctx, req.GetID(), req.GetOwner(), req.GetTitle(), req.GetText(), *startTime, *endTime)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &CreateEventResponse{
				Result: &CreateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, err
		}
	}
	protoEvent := &Event{
		ID:    event.ID,
		Title: event.Title,
		Text:  event.Text,
		Owner: event.Owner,
	}
	if protoEvent.StartTime, err = ptypes.TimestampProto(event.StartTime); err != nil {
		return nil, err
	}
	if protoEvent.EndTime, err = ptypes.TimestampProto(event.EndTime); err != nil {
		return nil, err
	}
	resp := &CreateEventResponse{
		Result: &CreateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

func (cs *CalendarServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*UpdateEventResponse, error) {

	var startTime *time.Time
	if req.GetStartTime() != nil {
		st, err := ptypes.Timestamp(req.GetStartTime())
		if err != nil {
			return nil, err
		}
		startTime = &st
	}

	var endTime *time.Time
	if req.GetEndTime() != nil {
		et, err := ptypes.Timestamp(req.GetEndTime())
		if err != nil {
			return nil, err
		}
		endTime = &et
	}

	event, err := cs.EventService.UpdateEvent(ctx, req.GetID(), req.GetOwner(), req.GetTitle(), req.GetText(), *startTime, *endTime)
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &UpdateEventResponse{
				Result: &UpdateEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, err
		}
	}
	protoEvent := &Event{
		ID:    event.ID,
		Title: event.Title,
		Text:  event.Text,
		Owner: event.Owner,
	}
	if protoEvent.StartTime, err = ptypes.TimestampProto(event.StartTime); err != nil {
		return nil, err
	}
	if protoEvent.EndTime, err = ptypes.TimestampProto(event.EndTime); err != nil {
		return nil, err
	}
	resp := &UpdateEventResponse{
		Result: &UpdateEventResponse_Event{
			Event: protoEvent,
		},
	}
	return resp, nil
}

func (cs *CalendarServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*DeleteEventResponse, error) {

	err := cs.EventService.DeleteEvent(ctx, req.GetID())
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &DeleteEventResponse{
				Result: &DeleteEventResponse_Error{
					Error: string(berr),
				},
			}
			return resp, nil
		} else {
			return nil, err
		}
	}

	resp := &DeleteEventResponse{
		Result: &DeleteEventResponse_Error{
			Error: "",
		},
	}
	return resp, nil
}

func (cs *CalendarServer) GetEventsByTime(ctx context.Context, req *GetEventsByTimeRequest) (*GetEventsByTimeResponse, error) {

	events, err := cs.EventService.GetEventsByTime(ctx, req.GetTimeType())
	if err != nil {
		if berr, ok := err.(errors.EventError); ok {
			resp := &GetEventsByTimeResponse{
				Error: string(berr),
			}
			return resp, nil
		} else {
			return nil, err
		}
	}
	var protoEvents []*Event
	for _, event := range events {
		protoEvent := &Event{
			ID:    event.ID,
			Title: event.Title,
			Text:  event.Text,
			Owner: event.Owner,
		}
		if protoEvent.StartTime, err = ptypes.TimestampProto(event.StartTime); err != nil {
			return nil, err
		}
		if protoEvent.EndTime, err = ptypes.TimestampProto(event.EndTime); err != nil {
			return nil, err
		}
		protoEvents = append(protoEvents, protoEvent)
	}

	resp := &GetEventsByTimeResponse{
		Event: protoEvents,
	}
	return resp, nil
}

func (cs *CalendarServer) Serve(addr string) error {
	// Create some standard server metrics.
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)
	grpc_prometheus.Register(grpcServer)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	RegisterCalendarServiceServer(s, cs)

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(grpcServer)

	return s.Serve(l)
}
