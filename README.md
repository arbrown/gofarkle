gofarkle
========

A tool for comparing AIs playing Farkle

[Farkle](http://en.wikipedia.org/wiki/Farkle) is a dice game with a number of variants.  This project aims to compare AIs playing a standardized version of Farkle (with support for adding rule variants in the future.)  Read [the basic rules](http://en.wikipedia.org/wiki/Farkle#Play) for an intro on how to play, then reference the scoring table below for  this project's standardized scoring rules.


| Dice                | Score | Dice                             | Score                                     |
| ----------          | ----  |--------------------------------- | ----------------------------------------  |
| **Each 1**          | 100   | **Four of a kind**               | 2x three of a kind <br/> for that number |
| **Each 5**          | 50    | **Five of a kind**               | 2x four of a kind <br/> for that number  |
| **Three 1's**       | 1000  | **Six of a kind**                | 2x five of a kind <br/> for that number  |
| **Three 2's**       | 200   | **Three Pairs**                  | 750                                       |
| **Three 3's**       | 300   | **Short Straigt (5)**            | 1500                                      |
| **Three 4's**       | 400   | **Long Straight (6)**            | 3000                                      |
| **Three 5's**       | 500   | **Three Farkles <br/> in a row** | -1000                                     |
| **Three 6's**       | 600   |
