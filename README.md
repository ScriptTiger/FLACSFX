[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://docs.google.com/forms/d/e/1FAIpQLSfBEe5B_zo69OBk19l3hzvBmz3cOV6ol1ufjh0ER1q3-xd2Rg/viewform)

# FLACSFX: FLAC SelF-eXtracting archive
FLACSFX is a minimal FLAC-to-WAV transcoder and multitrack FLAC mixer. The transcoder component transcodes embedded FLAC files to WAV files. The FLAC files are embedded by appending them to a precompiled FLACSFX executable. The multitrack FLAC mixer can mix the embedded FLAC tracks and output the mix either to a file or piped to standard output, to be played by VLC or other compatible applications without needing to write it to a file. This allows you to quickly and easily send losslessly compressed FLAC audio to someone who needs it as WAVs in order to reduce the file size in transit and increase transfer speeds, while also not requiring any technical know-how or additional software on the part of the recipient. And the mixer can be used either to extract the multitrack mix without having to store an additional mixed master track, or it can simply be used for quickly sampling the mixed audio when piped to something like VLC.

Usage: `flacsfx [options...]`
Argument                  | Description
--------------------------|-----------------------------------------------------------------------------------------------------
 `-flac`                  | Output FLAC
 `-mix`                   | Output mix
 `-o <directory\|file>`   | Destination directory, or file for mix
 `-b <16\|24\|32>`        | Bit depth of mix
 `-info`                  | Show stream info

`-` can be used in place of `<file>` to designate standard output as the destination for a mix.

Without any arguments, the embedded FLACs will be transcoded into the working directory within a new directory of the same name as the executable. The extracted audio files will either be titled using the pattern of `#_<executable name>`, or titled after the TITLE metadata tag if present within the FLAC. So, command-line usage is only optional and the end user can just execute the application as they would any other application for this default behavior.

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
