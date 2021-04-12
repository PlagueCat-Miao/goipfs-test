echo input 500k rand * X

read XX
mkdir output
cd output
seq $XX | xargs -i dd if=/dev/urandom of={}.dat bs=500k count=1 
