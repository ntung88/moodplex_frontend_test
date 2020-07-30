*Frontend Documentation*
======
Initial Design:-
<p align="center">
  <img src="https://github.com/sohamsankaran/moodplex/blob/master/frontend/content-img/homev1.png"><br>
  <b><ins>Initial Landing Page Design</ins></b>
</p>

Currently, the mood pages are not auto-generated (however, it is a fairly trivial fix) and "Happy" is the only working mood page. The user is subsequently directed to the corresponding mood page. 

<p align="center">
  <img src="https://github.com/sohamsankaran/moodplex/blob/master/frontend/content-img/ratev1.png"><br>
  <b><ins>Initial Rating Page Design</ins></b>
</p>

The rating page currently has a placeholder in place of content (from postgres). The user is able to rate whether they like the content or not with "No" or "Yes" buttons. Results are saved in json format as follows: item: 'response', dataType : 'json'. This can also be passed as a csv or another format that is preferable. The javascript for this is well commented so ot can be connected with the backend without confusion. 
