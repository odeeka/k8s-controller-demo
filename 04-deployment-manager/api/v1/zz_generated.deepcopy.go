// DeepCopy methods for AppDeployment types.
// AppDeploymentSpec and AppDeploymentStatus contain only value types (strings,
// int32) so no special map/slice handling is required — a plain struct copy works.
package v1

import runtime "k8s.io/apimachinery/pkg/runtime"

// ---------- AppDeployment ----------

func (in *AppDeployment) DeepCopyInto(out *AppDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *AppDeployment) DeepCopy() *AppDeployment {
	if in == nil {
		return nil
	}
	out := new(AppDeployment)
	in.DeepCopyInto(out)
	return out
}

func (in *AppDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// ---------- AppDeploymentList ----------

func (in *AppDeploymentList) DeepCopyInto(out *AppDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *AppDeploymentList) DeepCopy() *AppDeploymentList {
	if in == nil {
		return nil
	}
	out := new(AppDeploymentList)
	in.DeepCopyInto(out)
	return out
}

func (in *AppDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
