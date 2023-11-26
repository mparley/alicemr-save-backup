# Alice: Madness Returns Checkpoint Backupper

Got softlocked by the elevator bug in chapter 1 and was baffled that the game
doesn't keep any extra checkpoint saves or allows the user to save at all. So I
threw this together quickly. Just checks the last modified time of the checkpoint 
file every 5 minutes and if was updated it makes a copy. Keep it open as you play
and close when you are done. Probably not the best or most efficient way to do it 
but w/e. It works for me but use at your own risk.