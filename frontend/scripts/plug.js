document.getElementById("plug").innerHTML = `
<p style="line-height: 200%;">
        The Moodplex project is led by Soham Sankaran (Twitter:
        <a href="https://twitter.com/sohamsankaran" target="_blank">@sohamsankaran</a>), and
    is a production of <a href="https://honestyisbest.com/" target="_blank">Honesty Is Best</a>.
        <br>
        The Moodplex developers are: <a href="https://www.linkedin
        .com/in/nathan-tung-a3b562158/" target="_blank">Nathan Tung</a>,
        <a href="https://www.linkedin.com/in/kwilliamnrys/" target="_blank">Kyle
          Reyes</a>,
        and <a href="https://keethu-ram.github.io/" target="_blank">Keethu Ramalingam</a>. 
        <br>
        If you're interested in computer 
        science research, consider subscribing to the <a 
            href="https://honestyisbest.com/segfault" target="_blank">segfault podcast</a>. 
        <br>
        If you're a software engineer with 
        experience in compilers and building domain-specific languages, apply
    to work at <a href="https://pashi.com/" target="_blank">Pashi</a> by emailing
        soham [at] pashi [dot] com.
        </p>
<form action="https://buttondown.email/api/emails/embed-subscribe/moodplex" method="post" target="popupwindow" 
          onsubmit="window.open('https://buttondown.email/moodplex', 'popupwindow')" class="subscribe-form">
          <label for="subscribe-indianhistory-email" class="subscribe-label">&nbsp;Enter your email to subscribe to Moodplex updates</label>
          <input type="email" name="email" id="subscribe-indianhistory-email" placeholder="Type your email address here..." class="subscribe-textbox">
          <input type="hidden" value="1" name="embed">
          <input type="submit" value="Subscribe" class="subscribe-submit">
        </form>
`;