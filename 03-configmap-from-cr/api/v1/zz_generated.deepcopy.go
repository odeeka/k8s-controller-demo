// DeepCopy methods for ConfigSource types.
//
// ConfigSourceSpec contains a map[string]string (the Data field).
// Maps are reference types in Go, so we must deep copy them — otherwise two
// objects would share the same underlying map and modifying one would affect
// the other.
package v1

import runtime "k8s.io/apimachinery/pkg/runtime"

// ---------- ConfigSourceSpec ----------

// DeepCopyInto copies all fields, including the Data map.
func (in *ConfigSourceSpec) DeepCopyInto(out *ConfigSourceSpec) {
	*out = *in
	// Copy the map: allocate a new map and copy each entry.
	// We cannot just do `out.Data = in.Data` because that would make both
	// point to the same map in memory.
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// ---------- ConfigSource ----------

func (in *ConfigSource) DeepCopyInto(out *ConfigSource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec) // delegates to ConfigSourceSpec.DeepCopyInto above
	out.Status = in.Status
}

func (in *ConfigSource) DeepCopy() *ConfigSource {
	if in == nil {
		return nil
	}
	out := new(ConfigSource)
	in.DeepCopyInto(out)
	return out
}

func (in *ConfigSource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// ---------- ConfigSourceList ----------

func (in *ConfigSourceList) DeepCopyInto(out *ConfigSourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ConfigSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *ConfigSourceList) DeepCopy() *ConfigSourceList {
	if in == nil {
		return nil
	}
	out := new(ConfigSourceList)
	in.DeepCopyInto(out)
	return out
}

func (in *ConfigSourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
