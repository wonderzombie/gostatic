### gostatic

static site generator.

~~not working. broken.~~ work in progress.

### TODO

*  add support for self-contained templating (e.g. mustache)
    *  use case for this by itself? any reason not to go ahead with the next one?
*  add support for templates (e.g. _templates)
    *  need to spec out which data is available.
    *  tags would be a good way to exercise this.
*  add support for auto-linking golang.org/pkg for packages (?)
    *  stupid simple idea: search for tokens that look like pkg/os/exec and auto-link those.
        *  could do it before markdown processing step, where we s/[pkg/foo/bar]/[foo/bar][foo/bar]/ and then reference-based link at page bottom.
    *  another idea: custom tag. either manually parsed, or somehow made available via template language.
*  add support for table of contents
    *  should this happen at the blog generation level? yes for site-wide TOC, yes? 
    *  or at the templating language level? for a page TOC, yes. and some markdown variants already do this.