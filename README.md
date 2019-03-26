### deoxysii - Deoxys-II-256-128 for Go

This package provides a "from-the-paper" implementation of the
[Deoxys-II-256-128 v1.41][1] algorithm from the [final CAESAR portfolio][2].

#### Implementations

 * (`ct64`) Portable constant time implementation (Extremely Slow).

 * (`aesni`) SSSE3 + AESNI implementation for `amd64`

 * (`vartime`) Portable and variable time (insecure) implementation,
   for illustrative purposes (tested/benchmarked but never reachable
   or usable by external consumers).

#### Notes

Performance for the AES-NI implementation still has room for improvement,
however given that the Deoxys-BC-385 tweakable block cipher has 3 more
rounds than AES-256, and Deoxys-II will do two passes over the data
payload, it is likely reasonably close to what can be expected.

The pure software constant time implementation would benefit considerably
from vector optimizations as the amount of internal paralleism is quite
high, making it well suited to be implemented with [bitslicing][3].
Additionally a rather ludicrous amount of time is spent implementing the
`h` permutation in software, that can be replaced with a single `PSHUFB`
instruction.

[1]: https://competitions.cr.yp.to/round3/deoxysv141.pdf
[2]: https://competitions.cr.yp.to/caesar-submissions.html
[3]: https://eprint.iacr.org/2009/129.pdf