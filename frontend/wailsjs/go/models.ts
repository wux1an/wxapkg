export namespace main {
	
	export class FileFilter {
	    DisplayName: string;
	    Pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayName = source["DisplayName"];
	        this.Pattern = source["Pattern"];
	    }
	}

}

export namespace wechat {
	
	export class PathScanResult {
	    Paths: string[];
	    Logs: string;
	
	    static createFrom(source: any = {}) {
	        return new PathScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Paths = source["Paths"];
	        this.Logs = source["Logs"];
	    }
	}
	export class UnpackOptions {
	    EnableDecrypt: boolean;
	    EnableJsBeautify: boolean;
	    EnableHtmlBeautify: boolean;
	    EnableJsonBeautify: boolean;
	    OutputDir: string;
	    SavePath: string;
	
	    static createFrom(source: any = {}) {
	        return new UnpackOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EnableDecrypt = source["EnableDecrypt"];
	        this.EnableJsBeautify = source["EnableJsBeautify"];
	        this.EnableHtmlBeautify = source["EnableHtmlBeautify"];
	        this.EnableJsonBeautify = source["EnableJsonBeautify"];
	        this.OutputDir = source["OutputDir"];
	        this.SavePath = source["SavePath"];
	    }
	}
	export class WxapkgItem {
	    UUID: string;
	    WxId: string;
	    AppName: string;
	    AppNameSource: string;
	    Location: string;
	    EncryptKey: string;
	    Size: number;
	    IsDir: boolean;
	    LastModifyTime: number;
	    WxapkgFilePaths: string[];
	    UnpackStatus: string;
	    UnpackCurrent: number;
	    UnpackTotal: number;
	    UnpackProgress: number;
	    UnpackCurrentFile: string;
	    UnpackSavePath: string;
	    UnpackErrorMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new WxapkgItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.WxId = source["WxId"];
	        this.AppName = source["AppName"];
	        this.AppNameSource = source["AppNameSource"];
	        this.Location = source["Location"];
	        this.EncryptKey = source["EncryptKey"];
	        this.Size = source["Size"];
	        this.IsDir = source["IsDir"];
	        this.LastModifyTime = source["LastModifyTime"];
	        this.WxapkgFilePaths = source["WxapkgFilePaths"];
	        this.UnpackStatus = source["UnpackStatus"];
	        this.UnpackCurrent = source["UnpackCurrent"];
	        this.UnpackTotal = source["UnpackTotal"];
	        this.UnpackProgress = source["UnpackProgress"];
	        this.UnpackCurrentFile = source["UnpackCurrentFile"];
	        this.UnpackSavePath = source["UnpackSavePath"];
	        this.UnpackErrorMessage = source["UnpackErrorMessage"];
	    }
	}

}

