include ../Makefile.inc

SIMDIR=../src/
TESTMSG=To test the model launch \'Main.out\' in .$(OUTDIR) or the scripts in the main directory
MAINDEPS= $(SIMDIR)DT.6 $(SIMDIR)Sim.6 $(SIMDIR)Local.6 $(SIMDIR)Shared.6 $(SIMDIR)Random.6
ALLDEPS= Main.out

all: $(ALLDEPS)


Main.out:	Main.6 
	$(LD) -o ../$(OUTDIR)/Main.out Main.6
	$(ECHO) $(TESTMSG)

Main.6:	Main.go $(MAINDEPS)
	$(CC) Main.go

$(SIMDIR)Sim.6: force_look
	$(CD) $(SIMDIR); make Sim.6

$(SIMDIR)DT.6: force_look
	$(CD) $(SIMDIR); make DT.6

$(SIMDIR)Local.6: force_look
	$(CD) $(SIMDIR); make Local.6
	
$(SIMDIR)Shared.6: force_look
	$(CD) $(SIMDIR); make Shared.6
	
$(SIMDIR)Random.6: force_look
	$(CD) $(SIMDIR); make Random.6

clean:
	$(RM) *.8 *.6 *~

force_look:
	true
