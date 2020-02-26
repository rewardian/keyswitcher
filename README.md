# keyswitcher

**Linux is the only supported OS at the moment.**

Transform any Latin-script keyboard or numeric keypad into a macro keyboard (similar to the [Elgato StreamDeck](https://www.elgato.com/en/gaming/stream-deck)).

You can use keyswitcher to:
- Access global hotkeys configured within applications (e.g. [OBS](https://obsproject.com/), video games)
- Output ASCII on keypress (à la [AutoKey](https://github.com/autokey/autokey))
- Launch executables, one-liners, and bash scripts on keypress

By default, any keys without an assigned shortcut will have the "Super" and "Shift" keys applied to any other pressed key. For example, if you press "1" on an active keyboard, it will be interpreted as "Super" + "Shift" + "1".

In this way, you can assign global hotkeys for applications, accessing those hotkeys easily from the active keyboard while retaining normal use on the primary keyboard.

Secondly, you can assign shortcuts to any key on the activated keyboard (see: [keyswitcher.yaml](#keyswitcheryaml)). Currently you can set "text" and "program" shortcuts, to either output the saved text or to execute the saved script or application.

## Installation

Requires:
- Linux distribution
- Go 1.13
- Sudo privileges

The application is only compatible with Linux OSes currently and is built using Go. Root privileges are needed to directly access /dev/input, but any applications executed through keyswitcher launch as the normal user or sudo user.

```
git clone https://github.com/rewardian/keyswitcher.git

cd keyswitcher

## External dependencies
go get -u github.com/c-bata/go-prompt \
github.com/gvalkov/golang-evdev \
github.com/rewardian/keybd_event \
gopkg.in/yaml.v2

## Build the binary and copy it into a system-wide $PATH
go build; sudo mv ./keyswitcher /usr/local/bin/;

## Optional: Copy the demo YAML file for text and program shortcuts
cp -ar ./keyswitcher.yaml $HOME/keyswitcher.yaml
```

Once the binary is installed, you can launch the application using ```sudo keyswitcher```.

## Usage

### First Time
```
[lhammond@tantalus keyswitcher]$ sudo keyswitcher
[sudo] password for lhammond:

Shortcut file /home/lhammond/keyswitcher.yaml is loaded.

Options:
---
list		 Display available input devices
grab		 Gain exclusive control over input device
load		 Load a specific YAML configuration file [default: $HOME/keyswitcher.yaml]
release	         Release any activated keyboards
parse		 Print the shortcuts loaded in keyswitcher
log		 Toggle whether shortcuts, keypresses, or none are logged to stdout
exit		 Exit
keyswitcher>
```
I'd recommend disconnecting the keyboard you'd like to use as a macro keyboard. Then, type ```list```.
```
keyswitcher> list

Use the command 'grab ID' to take control of the device:

1 : AT Translated Set 2 keyboard - /dev/input/event4
```
This is *probably* the keyboard you want to continue using normally.

Plug in the keyboard you'd like to use and run ```list``` again:
```
keyswitcher> list

Use the command 'grab ID' to take control of the device:

1 : HOLDCHIP USB Gaming Keyboard - /dev/input/event15
2 : HOLDCHIP USB Gaming Keyboard - /dev/input/event18
3 : AT Translated Set 2 keyboard - /dev/input/event4
```
There's two things of note here.

1. Devices polled from /dev/input often have names other than what's marketed. For example, my Amazon Basics keyboard is named "CHICONY USB Keyboard"; my Azio HUE is named "HOLDCHIP USB Gaming Keyboard".

   Until we have profiles, you may be best off connecting the keyboard after running ```list``` so to identify which is which.

2. Most keyboards with special buttons (like ```Play/Pause Music```) present two separate devices. Some manufacturers label those devices as a 'consumer control' while others do not; you may need trial and error to ```grab``` the right device (or just grab both!).

Since I know ```1 : HOLDCHIP USB Gaming Keyboard - /dev/input/event15``` is the right one, I'm going to type ```grab 1```.

```
keyswitcher> grab 1
HOLDCHIP USB Gaming Keyboard - [active]>
```
The CLI's cursor has updated to the keyboard name and the keyboard is now marked active.

If you copied ```keyswitcher.yaml``` into your $HOME directory (during Installation), you can use a few example shortcuts now.

### keyswitcher.yaml
By default, ``$HOME/keyswitcher.yaml`` loads on program start.

You can load other YAML shortcut files by using ``load keyswitcher_example.yml`` or ``load /opt/file/another_example.yaml``.

The format of valid YAML should be relatively self-explanatory:
```
shortcuts:
  KEY_S:
    type: 'text'
    value: 'Lincoln Hammond\n281-330-8004\nlincoln@acornesque.net\n'
  KEY_M:
    type: 'program'
    value: 'if pgrep spotify > /dev/null; then wmctrl -a Spotify; else /var/lib/snapd/snap/bin/spotify ; fi'
  KEY_SPACE:
    type: 'program'
    value: 'gedit'
```
Under shortcuts, each list item is identified by keypress labels (found in [keymaps.go](https://github.com/rewardian/keyswitcher/blob/master/keymaps.go)) like 'KEY_SPACE' (Space) or 'KEY_KP0' (Numpad 0).

Under the key, we have two fields, ``Type`` and ``Value``.

For ``Type``, you can use 'text' or 'program'. If the ``Type`` is **text**, the ``Value`` will be output where your cursor is currently located.

If the ``Type`` is **program**, the ``Value`` will be executed as the normal user which ran ``sudo keyswitcher``.

Based on the ``keyswitcher.yaml`` file above, our `KEY_M` shortcut functions so that:
- If we press the 'M' key on an active keyboard,
- The bash one-liner ``if pgrep spotify > /dev/null; then wmctrl -a Spotify; else /var/lib/snapd/snap/bin/spotify ; fi`` will execute.

Edit ``$HOME/keyswitcher.yaml`` with your favorite text editor, or create a new YAML file in the format shown above, to add new shortcuts.

### Logging
Whenever a shortcut is executed, by default, its execution is logged to stdout in your terminal:
```
HOLDCHIP USB Gaming Keyboard - [active]> [KEY_S]: Outputting text: Lincoln Hammond\n281-330-8004\nlincoln@acornesque.net\n
[KEY_S]: Outputting text: Lincoln Hammond\n281-330-8004\nlincoln@acornesque.net\n
[KEY_M]: Executing program: if pgrep spotify > /dev/null; then wmctrl -a Spotify; else /var/lib/snapd/snap/bin/spotify ; fi
[KEY_M]: Executing program: if pgrep spotify > /dev/null; then wmctrl -a Spotify; else /var/lib/snapd/snap/bin/spotify ; fi
[KEY_SPACE]: Executing program: gedit
```

You can change the level of output by typing ``log``.
- **Verbose logging**: Outputs any keypress or shortcut execution, for debugging purposes or discovering which key label to use for a shortcut.
- **Normal logging**: Default. Output shortcut execution key, execution type, and commands or text.
- **Disabled**: Nothing is output to stdout from the active keyboards' use.

### Parse

You can use the command ``parse`` within keyswitcher to also output which shortcuts are currently loaded:
```
HOLDCHIP USB Gaming Keyboard - [active]> parse
Key: KEY_S
	Type: text
	Value: Lincoln Hammond\n281-330-8004\nlincoln@acornesque.net\n
Key: KEY_M
	Type: program
	Value: if pgrep spotify > /dev/null; then wmctrl -a Spotify; else /var/lib/snapd/snap/bin/spotify ; fi
Key: KEY_SPACE
	Type: program
	Value: gedit
```

### Release
To release any active keyboards, either exit the program (``exit`` or ``Ctrl+D``) or run ``release`` and press a key on the soon-to-be deactivated keyboard:
```
HOLDCHIP USB Gaming Keyboard - [active]> release
Press a key on the captured keyboard(s) to release.
keyswitcher>
```

## Contribute
If you have any questions, concerns, or feedback on the project, please feel free to:
- open an issue
- contact me at lincoln@acornesque.net
- or open a pull request
