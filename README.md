# TimedVoting

a small package for timed voting

## description
* Object `TimedVote` is used for timed vote , initial with parameters :
  * Vote : depands on IVote interface , you can set it up with defult function NewBasicVote , or create a new one with the interface .
  * Duration : set it for voting duration .
  * FeedBacks : depands on IVoteFeedBacks