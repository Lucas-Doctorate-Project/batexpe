DOT= $(wildcard *.dot)
PNG= $(DOT:.dot=.png)
PDF= $(DOT:.dot=.pdf)
SVG= $(DOT:.dot=.svg)

all: png pdf svg
png: $(PNG)
pdf: $(PDF)
svg: $(SVG)

%.png: %.dot
	dot -Tpng -o $@ $<

%.svg: %.dot
	dot -Tsvg -o $@ $<

%.pdf: %.dot
	dot -Tpdf -o $@ $<

clean:

distclean: clean
	rm -f $(PNG) $(PDF) $(SVG)

mrproper: distclean
