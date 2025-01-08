package handlers

const SystemPrompt = `You are a specialized content summarizer. You will receive structured web content that may contain some or all of the following components:
- Title: The main title of the content
- Headers: Important section headings with their types
- Meta Description: A brief overview of the content
- Main Content: The primary content body
- Paragraphs: Individual content sections
- Relevant Links: Related resources and references

Instructions for summarization:
1. Analyze available fields and build context from what is present, adapting if certain fields are empty
2. If title exists, use it to establish the main topic
3. If meta description is present, use it to support the overall context
4. For available headers, use them to understand the content structure
5. Combine information from main content and paragraphs to form a coherent narrative
6. Only reference links if they provide crucial context to the main topic

Key guidelines:
- Skip empty fields without mentioning their absence
- Connect information across available fields to build a complete picture
- Maintain context even with partial information
- Focus on creating a fluid, natural summary based on available content
- Use clear, direct language
- Keep the summary concise while preserving key information

Your goal is to deliver a coherent, well-structured summary regardless of which fields are present in the input.`

const SampleString = `Hollywood productions and events halted due to raging LA wildfires

From CNN's Elizabeth Wagmeister

The devastating wildfires that continue to ravage the celebrity-filled enclave of Pacific Palisades and other areas have forced a production shutdown across Los Angeles, as well as the cancelation of a number of key award season events that were set to take place this week.

The Critics Choice Awards, originally set to be held this Sunday in Santa Monica, have been postponed to January 26.

"This unfolding tragedy has already had a profound impact on our community. All our thoughts and prayers are with those battling the devastating fires and with all who have been affected," said Critics Choice Awards CEO Joey Berlin in a statement.

The Critics Choice Awards was set to be the second major televised Hollywood award show for the 2025 season, following last weekend's Golden Globes.

The award show was set to be held at The Barker Hanger, a venue in Santa Monica, not far from the Pacific Palisades where fire has destroyed at least 1,000 structures and burned more than 5,000 acres. Evacuation orders have also reached residents in Santa Monica where the award show was set to be held.

Amid the ongoing wildfires in Southern California, a number of glitzy Hollywood events and red carpet premieres have also been canceled.

The in-person nominations for the 31st annual Screen Actors Guild Awards were canceled on Wednesday morning, instead being announced via press release.

The annual AFI Awards luncheon, which was set to be held on January 10, will be rescheduled. And the annual BAFTA Tea Party, a key stop in the Oscars race set for January 11 at the Four Seasons Hotel in Beverly Hills, has been canceled, the organization announced.

Many Hollywood productions have been forced to stop filming as well, amid the high winds, smoke and dangerous fires.

More than a dozen shows that shoot in Los Angeles have halted production, according to The Hollywood Reporter, including "Grey's Anatomy," "Hacks," "Suits L.A.," "NCIS" and "The Price Is Right." Late night shows, like ABC's "Jimmy Kimmel Live!" and CBS' "After Midnight," will also cease production on Wednesday, per Variety, which reports that the situation will be monitored for Thursday's shows.

3 min ago

Another factor fueling the wildfires? A thirsty atmosphere

From CNN's Rachel Ramirez

A few major factors have contributed to the wildfires currently raging in Southern California, climate scientists say: strong winds, the lack of precipitation, and an increase in evaporative demand, also known as the "thirst of the atmosphere."

"It's basically how dry and thirsty the atmosphere is, and that correlates highly with wildfire potential," Michael Mann, a climate scientist at the University of Pennsylvania, told CNN's Brianna Keilar. "And we're seeing very high levels of evaporative demand precisely in those regions where these wildfires have broken out."

Warmer temperatures increase the amount of water the atmosphere can absorb, which then dries out the landscape. When the atmosphere sucks out the moisture from the soil without returning that water in the form of precipitation like rain, there's going to be less water available to those plants.

Daniel Swain, another climate scientist at the University of California, Los Angeles, said this thirstiness of the atmosphere has been "anomalously high" in the western parts of Los Angeles County, where the Palisades Fire is burning, as well as in the mountains where the Eaton Fire is spreading.

"Had we seen significant or widespread precipitation in the weeks and months leading up to this event, we would not be seeing the extent of devastation we are currently seeing," Swain said.`
