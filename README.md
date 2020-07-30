# moodplex
Code and docs for moodplex v1

Interacting with Backend:

- SQL Database structure {  
    id SERIAL,  
    rating float64, (elo score)  
    source TEXT, (url of specific image, like subdata)  
    category TEXT, (happy, sad, etc.)  
    agridataSource TEXT, (url of entire post)  
    website TEXT (twitter, reddit, hackernews, etc.)  
    }

** source and agridataSource are the same if the content is just text and the
 same url if the post/image are one and the same

- elo.go module path: "example.com/moodplex/backend"

- to run elo.go: "go run elo.go"

Handlers:

1) Post to "/match" with the result of a match (whenever the user inputs their 
preference given the two posts) to have it recorded in the database and elo 
scores updated
    - will need the two post id's and their scores in the matchup
2) "/posts" to get a single random post from the database with the given mood 
category