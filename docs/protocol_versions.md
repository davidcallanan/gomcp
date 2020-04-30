# Java Edition Protocol Versions

The protocol versions in this package differ from those encoded into the protocol. This is to remove the discrepency between pre-netty and post-netty version numbers. The numbers are also represented with 4 hex characters to save space.

## Pre-Netty Version Numbers

This refers to versions before or equal to the `13w39b` 1.7 snapshot.

The encoded protocol version numbers can be found [here](https://wiki.vg/Protocol_version_numbers#Versions_before_the_Netty_rewrite).

The version number used in this package for each of these versions can be obtained by converting its encoded protocol version number to 4-digit hex.

## Post-Netty Version Numbers

This refers to versions after or equal to the `13w41a` 1.7 snapshot.

The encoded protocol version numbers can be found [here](https://wiki.vg/Protocol_version_numbers#Versions_after_the_Netty_rewrite).

The version number used in this package for each of these versions can be obtained by adding decimal `81` to its encoded protocol version number and then converting it to 4-digit hex.

## Other Minecraft Editions

To differentiate between protocol versions in Java Edition and other editions of Minecraft, the Java Edition version number may be prefixed with a `j` character.
