# ggit: git in golang

* ggit is an (aspiring) implementation of Git written in Go.
* ggit is a single-executable Git that compiles on Go-supported systems, including Windows.
* ggit is a library for programmatically interacting with Git repositories.

# Status

ggit is alpha. We currently support reading operations and an API around the following constructs: 
blobs, trees, tags, commits, refs, packed refs, git packs, and index files. ggit can parse revisions
and has a command suite similar to git.

# Changelog

* **0.1.0** (TBA) Minimal viable reading API with test suite and benchmarks.
* **0.0.2** Pack parsing and testing framework.
* **0.0.1** Licensing and API improvements.
* **0.0.0** Basic objects and command suite.

# Copyright

    Copyright (c) 2012 The ggit Authors

# Authors

This project is authored and maintained by:

    Jake Brukhman <brukhman@gmail.com>
    Michael Bosworth <michael.a.bosworth@gmail.com>

# License

<a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/3.0/deed.en_US"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by-nc-nd/3.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/3.0/deed.en_US">Creative Commons Attribution-NonCommercial-NoDerivs 3.0 Unported License</a>.

Under this license you are free, without having to seek permission, to use ggit and to view its source in accordance with the license terms. For permissions beyond the scope of this license, please contact the authors as specified in the previous section.

You must keep intact the copyright notice above and give attribution to the authors of this work. If you would like to become a co-author or contributor to ggit, please contact us.