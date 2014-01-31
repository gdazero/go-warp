include ../Makefile.inc

ALLDEPS= Random.6 Const.6 DT.6 Heap.6 Communication.6 Gvt.6 Sim.6 Local.6 Shared.6

all: $(ALLDEPS)


Heap.6:	Heap.go DT.6 Const.6
	$(CC) Heap.go

Random.6:	Random.go
	$(CC) Random.go

Sim.6:	Sim.go DT.6 Communication.6 Local.6 Const.6 Gvt.6 Shared.6
	$(CC) Sim.go

DT.6:	DT.go Const.6
	$(CC) DT.go

Const.6:	Const.go
	$(CC) Const.go

Communication.6:	Communication.go DT.6 Const.6
	$(CC) Communication.go

Local.6:	Local.go DT.6 Heap.6 Const.6
	$(CC) Local.go

Gvt.6:	Gvt.go Const.6 DT.6 Shared.6
	$(CC) Gvt.go

Shared.6:	Shared.go DT.6 Local.6
	$(CC) Shared.go

clean:
	$(RM) *.6 *~
