package deltautil

import (
	"fmt"

	pb "github.com/earthly/cloud-api/logstream"
	"google.golang.org/protobuf/proto"
)

// Version is the version of the deltautil package.
const Version = 2

// WithDeltaManifest takes a delta and and a manifest and returns the result of
// applying the delta to the manifest. The original passed-in manifest is not
// changed. If the delta would not have any effect on the manifest, the original
// manifest is returned.
func WithDeltaManifest(m *pb.RunManifest, d *pb.Delta) (*pb.RunManifest, error) {
	if m.GetVersion() != 0 && m.GetVersion() != Version {
		return nil, fmt.Errorf("unsupported manifest version %d", m.GetVersion())
	}
	if d.GetVersion() != 0 && d.GetVersion() != Version {
		return nil, fmt.Errorf("unsupported delta version %d", d.GetVersion())
	}
	var dm *pb.DeltaManifest
	switch d.GetDeltaTypeOneof().(type) {
	case *pb.Delta_DeltaManifest:
		dm = d.GetDeltaManifest()
	default:
		// No action needed if this is not a manifest delta.
		return m, nil
	}
	var ret *pb.RunManifest
	switch dm.GetDeltaManifestOneof().(type) {
	case *pb.DeltaManifest_ResetAll:
		ret = proto.Clone(dm.GetResetAll()).(*pb.RunManifest)
	case *pb.DeltaManifest_Fields:
		ret = proto.Clone(m).(*pb.RunManifest)
		f := dm.GetFields()
		if f.GetStartedAtUnixNanos() != 0 {
			ret.StartedAtUnixNanos = f.GetStartedAtUnixNanos()
		}
		if f.GetEndedAtUnixNanos() != 0 {
			ret.EndedAtUnixNanos = f.GetEndedAtUnixNanos()
		}
		if f.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
			ret.Status = f.GetStatus()
		}
		if f.GetHasFailure() {
			ret.Failure = f.GetFailure()
		}
		if f.GetMainTargetId() != "" {
			ret.MainTargetId = f.GetMainTargetId()
		}
		for targetID, t2 := range f.GetTargets() {
			if ret.Targets == nil {
				ret.Targets = make(map[string]*pb.TargetManifest)
			}
			t, ok := ret.Targets[targetID]
			if !ok {
				t = &pb.TargetManifest{}
				ret.Targets[targetID] = t
			}

			if t2.GetName() != "" {
				t.Name = t2.GetName()
			}
			if t2.GetCanonicalName() != "" {
				t.CanonicalName = t2.GetCanonicalName()
			}
			if len(t2.GetOverrideArgs()) > 0 {
				t.OverrideArgs = append([]string{}, t2.GetOverrideArgs()...)
			}
			if t2.GetInitialPlatform() != "" {
				t.InitialPlatform = t2.GetInitialPlatform()
			}
			if t2.GetFinalPlatform() != "" {
				t.FinalPlatform = t2.GetFinalPlatform()
			}
			if t2.GetRunner() != "" {
				t.Runner = t2.GetRunner()
			}
			if t2.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
				t.Status = t2.GetStatus()
			}
			if t2.GetStartedAtUnixNanos() != 0 {
				t.StartedAtUnixNanos = t2.GetStartedAtUnixNanos()
			}
			if t2.GetEndedAtUnixNanos() != 0 {
				t.EndedAtUnixNanos = t2.GetEndedAtUnixNanos()
			}
		}
		for commandID, c2 := range f.GetCommands() {
			if ret.Commands == nil {
				ret.Commands = make(map[string]*pb.CommandManifest)
			}
			c, ok := ret.Commands[commandID]
			if !ok {
				c = &pb.CommandManifest{}
				ret.Commands[commandID] = c
			}

			if c2.GetName() != "" {
				c.Name = c2.GetName()
			}
			if c2.GetTargetId() != "" {
				c.TargetId = c2.GetTargetId()
			}
			if c2.GetCategory() != "" {
				c.Category = c2.GetCategory()
			}
			if c2.GetPlatform() != "" {
				c.Platform = c2.GetPlatform()
			}
			if c2.GetStatus() != pb.RunStatus_RUN_STATUS_UNKNOWN {
				c.Status = c2.GetStatus()
			}
			if c2.GetHasCached() {
				c.IsCached = c2.GetIsCached()
			}
			if c2.GetHasLocal() {
				c.IsLocal = c2.GetIsLocal()
			}
			if c2.GetHasInteractive() {
				c.IsInteractive = c2.GetIsInteractive()
			}
			if c2.GetStartedAtUnixNanos() != 0 {
				c.StartedAtUnixNanos = c2.GetStartedAtUnixNanos()
			}
			if c2.GetEndedAtUnixNanos() != 0 {
				c.EndedAtUnixNanos = c2.GetEndedAtUnixNanos()
			}
			if c2.GetHasHasProgress() {
				c.HasProgress = c2.GetHasProgress()
			}
			if c2.GetHasProgress() {
				c.Progress = c2.GetProgress()
			}
			if c2.GetErrorMessage() != "" {
				c.ErrorMessage = c2.GetErrorMessage()
			}
			if c2.GetHasSourceLocation() {
				c.SourceLocation = c2.GetSourceLocation()
			}
		}
	}
	return ret, nil
}
