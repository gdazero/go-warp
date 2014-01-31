include Makefile.inc

ALLDEPS=$(SRCDIR)/Makefile $(MODELDIR)/Makefile

all :   $(ALLDEPS)
	$(ECHO)
	$(ECHO) $(SUCCMSG)

$(SRCDIR)/Makefile : force_look
	$(ECHO) looking into $(SRCDIR):
	$(CD) $(SRCDIR); $(MAKEALL)

$(MODELDIR)/Makefile : force_look
	$(ECHO)
	$(ECHO) looking into $(MODELDIR)
	$(CD) $(MODELDIR); $(MAKEALL)


clean:	
	$(ECHO) $(CLEANMSG)
	$(RM) $(SRCDIR)*.6 $(SRCDIR)*.8
	$(RM) $(OUTDIR)*.out
	$(RM) $(MODELDIR)*.6 $(MODELDIR)*.8

cleanlog:
	$(ECHO) $(LOGMSG)
	$(RM) $(LOGDIR)*

force_look:
	true
