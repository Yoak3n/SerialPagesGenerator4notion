export namespace api {
	
	export class Section {
	    aid: number;
	    long_title: string;
	    short_title: string;
	    cover: string;
	
	    static createFrom(source: any = {}) {
	        return new Section(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.aid = source["aid"];
	        this.long_title = source["long_title"];
	        this.short_title = source["short_title"];
	        this.cover = source["cover"];
	    }
	}
	export class Bangumi {
	    season_id: string;
	    title: string;
	    bangumi_cover: string;
	    main: Section[];
	    sub: Section[];
	
	    static createFrom(source: any = {}) {
	        return new Bangumi(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.season_id = source["season_id"];
	        this.title = source["title"];
	        this.bangumi_cover = source["bangumi_cover"];
	        this.main = this.convertValues(source["main"], Section);
	        this.sub = this.convertValues(source["sub"], Section);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class VideoInfo {
	    name: string;
	    aid: number;
	    bvid: string;
	    titles: string[];
	    cover: string;
	    video_view: number;
	    video_time: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.aid = source["aid"];
	        this.bvid = source["bvid"];
	        this.titles = source["titles"];
	        this.cover = source["cover"];
	        this.video_view = source["video_view"];
	        this.video_time = source["video_time"];
	    }
	}

}

export namespace config {
	
	export class Config {
	    database_id: string;
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.database_id = source["database_id"];
	        this.token = source["token"];
	    }
	}

}

export namespace model {
	
	export class Name {
	    // Go type: struct { Name string "json:\"name\"" }
	    select: any;
	
	    static createFrom(source: any = {}) {
	        return new Name(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.select = this.convertValues(source["select"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Title {
	    type: string;
	    // Go type: struct { Content string "json:\"content\"" }
	    text: any;
	
	    static createFrom(source: any = {}) {
	        return new Title(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.text = this.convertValues(source["text"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EpisodeName {
	    title: Title[];
	
	    static createFrom(source: any = {}) {
	        return new EpisodeName(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = this.convertValues(source["title"], Title);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Episode {
	    number: number;
	
	    static createFrom(source: any = {}) {
	        return new Episode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.number = source["number"];
	    }
	}
	export class Data {
	    // Go type: Episode
	    Episode: any;
	    // Go type: EpisodeName
	    EpisodeName: any;
	    // Go type: Name
	    Name: any;
	    type: string;
	    database_id: string;
	
	    static createFrom(source: any = {}) {
	        return new Data(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Episode = this.convertValues(source["Episode"], null);
	        this.EpisodeName = this.convertValues(source["EpisodeName"], null);
	        this.Name = this.convertValues(source["Name"], null);
	        this.type = source["type"];
	        this.database_id = source["database_id"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

