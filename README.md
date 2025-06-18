# wingman - Superbird manager

Open source management tool for the Spotify Car Thing.

## Features

- Manage Android A/B metadata from misc partition (open source `phb` implementation)
  - Get A/B data
  - Set boot result
  - Set active slot
  - Reset and switch to slot A
  - Output as JSON for parsing in another program

## Usage

```
USAGE:
   wingman [global options] command [command options]

COMMANDS:
   ab       set a/b data
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

```
NAME:
   wingman ab - set a/b data

USAGE:
   wingman ab [command options]

OPTIONS:
   --json               output data in JSON format (default: false)
   --boot-result value  set the boot result. 0 for failure, 1 for success
   --slot value         set the active boot slot. 0 for A, 1 for B
   --reset              reset all boot data and switch back to slot A (default: false)
   --help, -h           show help
```

## Donate

Nocturne is a massive endeavor, and the team have spent everyday over the last few months making it a reality out of our passion for creating something that people like you love to use.

All donations are split between the four members of the Nocturne team, and go towards the development of future features. We are so grateful for your support!

[Buy Me a Coffee](https://buymeacoffee.com/brandonsaldan) | [Ko-Fi](https://ko-fi.com/brandonsaldan)

## Credits

This software was made possible only through the following individuals and open source programs:

- [Dominic Frye](https://github.com/itsnebulalol)

<br />

- [spsgsb/uboot](https://github.com/spsgsb/uboot/blob/buildroot-openlinux-201904-g12a/common/avb.c)

## License

This project is licensed under the **MIT** license.

---

> © 2025 Nocturne.

> "Spotify" and "Car Thing" are trademarks of Spotify AB. This software is not affiliated with or endorsed by Spotify AB.

> [usenocturne.com](https://usenocturne.com) &nbsp;&middot;&nbsp;
> GitHub [@usenocturne](https://github.com/usenocturne)
