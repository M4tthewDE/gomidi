document.addEventListener('DOMContentLoaded', function() {
	const VF = Vex.Flow;
	const canvas = document.getElementById('musicCanvas');
	const renderer = new VF.Renderer(canvas, VF.Renderer.Backends.CANVAS);

	const context = renderer.getContext();
	context.setFont('Arial', 10, '').setBackgroundFillStyle('#eed');

	const stave = new VF.Stave(10, 40, 500);
	stave.addClef('treble').addTimeSignature('4/4');
	stave.setContext(context).draw();
});
