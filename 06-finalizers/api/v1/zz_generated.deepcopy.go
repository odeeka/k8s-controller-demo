// DeepCopy methods for TrackedResource types.
// All fields are value types (strings) — straightforward copies.
package v1

import runtime "k8s.io/apimachinery/pkg/runtime"

// ---------- TrackedResource ----------

func (in *TrackedResource) DeepCopyInto(out *TrackedResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *TrackedResource) DeepCopy() *TrackedResource {
	if in == nil {
		return nil
	}
	out := new(TrackedResource)
	in.DeepCopyInto(out)
	return out
}

func (in *TrackedResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// ---------- TrackedResourceList ----------

func (in *TrackedResourceList) DeepCopyInto(out *TrackedResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TrackedResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *TrackedResourceList) DeepCopy() *TrackedResourceList {
	if in == nil {
		return nil
	}
	out := new(TrackedResourceList)
	in.DeepCopyInto(out)
	return out
}

func (in *TrackedResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
