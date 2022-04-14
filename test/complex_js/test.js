if (end == 0)
	deform = weighted ? test ? a : Utils.newFloatArray(deformLength) : Utils.newFloatArray(deformLength);
else {
	deform = Utils.newFloatArray(deformLength);
}