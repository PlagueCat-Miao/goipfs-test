num = 34

buildrandom:
	echo input 500k * X
	rm -rf output
	mkdir output
	cd output ;\
	seq "$(num)" | xargs -i dd if=/dev/urandom of={}.dat bs=500k count=1
clear:
	rm -rf output/*.dat
	rm ./*.dat
	rm output

