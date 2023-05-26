# wxapkg

**Install:**

Download the latest [release](https://github.com/wux1an/wxapkg/releases/latest) or build with `go install` by yourself

```
go install github.com/wux1an/wxapkg@latest
```

**Usage:**

![demo](demo.gif)

- scan mini program

    ```bash
    # wxapkg.exe scan --help
    wxapkg.exe scan
    ```

- unpack mini program

    ```bash
    # wxapkg.exe unpack --help     
    wxapkg.exe unpack -o out-dir -r "%USERPROFILE%\Documents\WeChat Files\Applet\wx00000000000000"
    ```

**Todo:**

- [x] scan more information (need networking）
- [x] json beautify
- [ ] javascript beautify
- [ ] html beautify
- [ ] auto export uri in files

**References:**

- decrypt: https://github.com/BlackTrace/pc_wxapkg_decrypt
- unpack: [https://gist.github.com/Integ/bcac5c21de5ea3...](https://gist.github.com/Integ/bcac5c21de5ea35b63b3db2c725f07ad)
- introduce: [https://misakikata.github.io/2021/03/%E5%BE%...](https://misakikata.github.io/2021/03/%E5%BE%AE%E4%BF%A1%E5%B0%8F%E7%A8%8B%E5%BA%8F%E8%A7%A3%E5%8C%85/)
- terminal ui: https://github.com/charmbracelet/bubbletea