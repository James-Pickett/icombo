# icombo

Declarative Image Combination

## Why?

I have been dabbling in board game design as a hobby with my family. The most toilsome part of the process for me is editing images for all the little icons. For example, let's say I need a set of icons to show that you can covert 2 logs into a board. I would need to open an image editor paste the various parts in there, make sure sizing is correct, etc. Then if I wanted to change the conversion to 3 logs into a board, I would have to open the image editor again. So I developed icombo to speed up this process and provide a declarative format (toml) that will handle joining the images, rotation and sizing to provide consistent results.

Now my workflow looks like this:
* go to [game-icons.net](https://game-icons.net/) and find the parts for the image I want to create

    <img src="./example/image_parts/log.png" width="50">
    <img src="./example/image_parts/arrow.png" width="50">
    <img src="./example/image_parts/board.png" width="50">
    
* drop them into my image_parts directory
* add a few lines to my [icombo.toml](./example/icombo.toml) file
    ```
    [[images]]
    name = "logs_to_boards"

        [[images.image_parts]]
        file_name = "log"
        count = 2

        # arrow faces up by default and rotates counter clockwise
        [[images.image_parts]]
        file_name = "arrow"
        rotation_degrees = 270

        [[images.image_parts]]
        file_name = "board"
    ```
* run icombo to produce logs_to_boards.png

    <img src="./example/output_images/logs_to_boards.png" width="200">

### Why TOML?
While my personal preference would by yml, I think toml is easer for non coder types

## Limitations
* only handles .png files
* only builds images horizontally from left to right
* cannot configure individual images with size
* no support for setting background color
* transparency untested

## Possible Improvments
* address limitations above
* instead of reading each image part from disk everytime, store them in memory
  - this may or may not be helpful, need to test, possibly trading memory for speed
* add flags for config file path, now it must be at `./icombo.toml`
* add functionality for handling other config formats (yml, json, etc)

## Credits
* [game-icons.net](https://game-icons.net/) for the art
* @ozonru for the [imaging package](https://github.com/disintegration/imaging) that made the image manipulation super easy
* @spf13 for the [viper package's](https://github.com/spf13/viper) awesome config tooling
