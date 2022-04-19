
var a = class extends CurveTimeline {
	constructor(frameCount, bezierCount, slotIndex) {
	  super(frameCount, bezierCount, [
		Property.rgb + "|" + slotIndex,
		Property.alpha + "|" + slotIndex
	  ]);
	}
}