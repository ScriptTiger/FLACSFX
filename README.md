[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# FLACSFX: FLAC SelF-eXtracting archive
FLACSFX is a minimal FLAC-to-WAV transcoder and multitrack FLAC mixer for embedded FLAC archives. The FLAC files are embedded and archived by appending them to a precompiled FLACSFX executable and can be either unindexed or indexed, for better performance when working with larger archives. The transcoder component transcodes embedded FLAC files either to WAV files or standard output, while the multitrack FLAC mixer can mix any combination of compatible embedded FLAC tracks and output the mix either to a file or standard output, allowing you to play source tracks or new mixes by piping them to FFmpeg, FFplay, VLC, or other compatible applications without needing to write a new file.

FLACSFX came about from my own experience in the content creation space and has slowly evolved to what it is today. As content creation is often a collaborative effort, FLACSFX has evolved to foster collaboration, independent of whatever video editors, DAWs, audio editors, and even operating systems team members might be working with. By adding new tracks and/or layers to an archive, it allows team members to experiment much more quickly without having to worry about compatibility issues. You can quickly and easily send off multiple tracks with varying versions of various tracks and allow other team members to experiment with different combinations and give their opinion, rather than having to render a separate version for each possible combination. You could also just use it to simply send off your finished audio to a video editor who's more than likely using a video editor application which doesn't support FLAC, as most don't, so you can get all the benefits of sending a FLAC and they can get all the benefits of receiving a WAV.

Usage: `flacsfx [options...]`
Argument                  | Description
--------------------------|-----------------------------------------------------------------------------------------------------
 `-i <#,#-#,...>`         | Index entries to include
 `-o <directory\|file>`   | Destination directory, or file for single entry, mix, or index
 `-flac`                  | Output FLAC, cannot be used with -mix
 `-mix`                   | Output mix to WAV
 `-bits <16\|24\|32>`     | Bit depth of mix
 `-attenuate`             | Attenuate linearly to prevent clipping in mix, dividing by number of tracks
 `-info`                  | Show stream info
 `-index`                 | Save the index to a file
 `-ignoreindex`           | Ignore the index

By default, a new directory is created to extract the audio files to. However, if only one entry is selected, or a mix is being created, the output is a file and not a directory. `-` can thus be used in place of `<file>` to designate standard output as the destination for a single entry or mix, allowing single entries and mixes to be sampled by compatible software, such as FFmpeg, FFplay, VLC, without needing to write a file.

Without any arguments, the embedded FLACs will be transcoded into the working directory within a new directory of the same name as the executable, or into the working directory directly if only one embedded FLAC is present. The extracted audio files will either be titled using the pattern of `<track number>_<executable name>`, or simply `<executable name>` if only one embedded FLAC is present, or titled after the TITLE metadata tag if present within the FLAC. So, command-line usage is only optional and the end user can just execute the application as they would any other application for this default behavior.

When you append FLACs to a FLACSFX exectuable, they are initially not indexed within the archive and additional time must be spent to build a running index. This time may not be noticeable for smaller archives, but increases as the size of the archive increases. For additional performance improvements, you can save the index to a file and append that file to the FLACSFX in the same way you would append FLAC files. And since the index is added to the end of the file, you must thus generate a new index each time new FLACs are appended, to ensure the new index is the last thing appended to the file. However, you may continue to add as many new FLACs and new indices as you want without issue. The size of the index itself is only a few bytes and is thus negligible overhead. You can also append old indices after you've appended new FLACs and this will allow you to index the new archive faster by reusing the old index for the old tracks and instantly jumping straight to indexing the new tracks.

You can also append FLACSFX executables to the front of existing archives and it will instantly convert that archive to whatever the target operating system is. As the FLACSFX executable itself is rather minimal, less than 2 megabytes, it's not much additional overhead. However, just to reiterate, when appending a FLACSFX executable to a pre-existing archive, the FLACSFX executable for the target operating system must be the first item listed when concatenating the files. If the pre-existing archive was indexed, you will also need to use the `-ignoreindex` argument to ensure the old index is not used, since the old index will be incorrect now that the new FLACSFX executable has been added and has changed the offsets of the contents.

For additional notes on the mixer package used, please refer to its documentation:  
https://github.com/ScriptTiger/mixerInG

# Appending FLAC files to a FLACSFX executable
Download the latest precompiled releases for the intended target system:  
https://github.com/ScriptTiger/FLACSFX/releases/latest

For appending FLAC files to a FLACSFX executable, issue one of the following commands.

For Windows:
```
copy /b "FLACSFX.exe"+"file1.flac"+"file2.flac"+... "MyFLACSFX.exe"
```

For Linux and Mac:
```
cat "FLACSFX" "file1.flac" "file2.flac" ... > "MyFLACSFX"
```

# Single-track FLACSFX
If you don't need any of the more advanced features of the newer version, such as built-in mixing and multitrack capabilities for embedding multiple audio files, the older version for simply embedding a single audio file is still available, although no longer maintained.

README:  
https://github.com/ScriptTiger/FLACSFX/tree/4ffca93e97a915359c7511f26735ca35a7eac0bd

Precompiled executables:  
https://github.com/ScriptTiger/FLACSFX/releases/tag/25AUG2024

# More About ScriptTiger

For more ScriptTiger scripts and goodies, check out ScriptTiger's GitHub Pages website:  
https://scripttiger.github.io/

[![Donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=MZ4FH4G5XHGZ4)
