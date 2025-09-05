export namespace lib {
	
	export class TruthTableVariable {
	    Name: string;
	    Value: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TruthTableVariable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Value = source["Value"];
	    }
	}
	export class TruthTableEntry {
	    Result: boolean;
	    Variables: TruthTableVariable[];
	
	    static createFrom(source: any = {}) {
	        return new TruthTableEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Result = source["Result"];
	        this.Variables = this.convertValues(source["Variables"], TruthTableVariable);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

