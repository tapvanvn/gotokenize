function setFrame(frame, time, value) {
	frame <<= 1;
	this.frames[frame] = time;
	this.frames[frame + 1] = value;
  }