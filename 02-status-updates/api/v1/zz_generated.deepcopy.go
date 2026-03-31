// DeepCopy methods for the Greeter types.
// See 01-minimal-controller for a detailed explanation of why these are needed.
package v1

import runtime "k8s.io/apimachinery/pkg/runtime"

// ---------- Greeter ----------

func (in *Greeter) DeepCopyInto(out *Greeter) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	// GreeterStatus contains metav1.Time (which wraps time.Time — a value type).
	// A plain struct copy is safe here; no pointers that need independent allocation.
	out.Status = in.Status
}

func (in *Greeter) DeepCopy() *Greeter {
	if in == nil {
		return nil
	}
	out := new(Greeter)
	in.DeepCopyInto(out)
	return out
}

func (in *Greeter) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// ---------- GreeterList ----------

func (in *GreeterList) DeepCopyInto(out *GreeterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Greeter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *GreeterList) DeepCopy() *GreeterList {
	if in == nil {
		return nil
	}
	out := new(GreeterList)
	in.DeepCopyInto(out)
	return out
}

func (in *GreeterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
