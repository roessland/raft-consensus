# raft-consensus

Implementation of Raft distributed consensus for learning purposes.

Based on the lectures by Martin Kleppmann:

- https://www.youtube.com/watch?v=uXEYuDwm7e4
- https://www.cl.cam.ac.uk/teaching/2122/ConcDisSys/dist-sys-notes.pdf

The goal is to serve a simple but resilient key-value store over HTTP.
It should work with multiple processes on localhost, and on multiple computers in a cluster.
It should be readable and avoid mutexes when possible.

## Backlog

- [ ] Improve test suite by adding fault injection and message reordering in HTTP proxy
- [ ] Verify that network transport used is compliant with Raft paper with respect to reordering.
- [ ] Implement broadcast using various message orderings
- [ ] Finish key/value database API
- [ ] Implement hashicorp/raft interface and try to run their test suite
