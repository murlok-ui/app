// +build !wasm

package app

// Code generated by go generate; DO NOT EDIT.

const (
	wasmExecJS = "// Copyright 2018 The Go Authors. All rights reserved.\n// Use of this source code is governed by a BSD-style\n// license that can be found in the LICENSE file.\n\n(() => {\n\t// Map multiple JavaScript environments to a single common API,\n\t// preferring web standards over Node.js API.\n\t//\n\t// Environments considered:\n\t// - Browsers\n\t// - Node.js\n\t// - Electron\n\t// - Parcel\n\n\tif (typeof global !== \"undefined\") {\n\t\t// global already exists\n\t} else if (typeof window !== \"undefined\") {\n\t\twindow.global = window;\n\t} else if (typeof self !== \"undefined\") {\n\t\tself.global = self;\n\t} else {\n\t\tthrow new Error(\"cannot export Go (neither global, window nor self is defined)\");\n\t}\n\n\tif (!global.require && typeof require !== \"undefined\") {\n\t\tglobal.require = require;\n\t}\n\n\tif (!global.fs && global.require) {\n\t\tglobal.fs = require(\"fs\");\n\t}\n\n\tif (!global.fs) {\n\t\tlet outputBuf = \"\";\n\t\tglobal.fs = {\n\t\t\tconstants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused\n\t\t\twriteSync(fd, buf) {\n\t\t\t\toutputBuf += decoder.decode(buf);\n\t\t\t\tconst nl = outputBuf.lastIndexOf(\"\\n\");\n\t\t\t\tif (nl != -1) {\n\t\t\t\t\tconsole.log(outputBuf.substr(0, nl));\n\t\t\t\t\toutputBuf = outputBuf.substr(nl + 1);\n\t\t\t\t}\n\t\t\t\treturn buf.length;\n\t\t\t},\n\t\t\twrite(fd, buf, offset, length, position, callback) {\n\t\t\t\tif (offset !== 0 || length !== buf.length || position !== null) {\n\t\t\t\t\tthrow new Error(\"not implemented\");\n\t\t\t\t}\n\t\t\t\tconst n = this.writeSync(fd, buf);\n\t\t\t\tcallback(null, n);\n\t\t\t},\n\t\t\topen(path, flags, mode, callback) {\n\t\t\t\tconst err = new Error(\"not implemented\");\n\t\t\t\terr.code = \"ENOSYS\";\n\t\t\t\tcallback(err);\n\t\t\t},\n\t\t\tread(fd, buffer, offset, length, position, callback) {\n\t\t\t\tconst err = new Error(\"not implemented\");\n\t\t\t\terr.code = \"ENOSYS\";\n\t\t\t\tcallback(err);\n\t\t\t},\n\t\t\tfsync(fd, callback) {\n\t\t\t\tcallback(null);\n\t\t\t},\n\t\t};\n\t}\n\n\tif (!global.crypto) {\n\t\tconst nodeCrypto = require(\"crypto\");\n\t\tglobal.crypto = {\n\t\t\tgetRandomValues(b) {\n\t\t\t\tnodeCrypto.randomFillSync(b);\n\t\t\t},\n\t\t};\n\t}\n\n\tif (!global.performance) {\n\t\tglobal.performance = {\n\t\t\tnow() {\n\t\t\t\tconst [sec, nsec] = process.hrtime();\n\t\t\t\treturn sec * 1000 + nsec / 1000000;\n\t\t\t},\n\t\t};\n\t}\n\n\tif (!global.TextEncoder) {\n\t\tglobal.TextEncoder = require(\"util\").TextEncoder;\n\t}\n\n\tif (!global.TextDecoder) {\n\t\tglobal.TextDecoder = require(\"util\").TextDecoder;\n\t}\n\n\t// End of polyfills for common API.\n\n\tconst encoder = new TextEncoder(\"utf-8\");\n\tconst decoder = new TextDecoder(\"utf-8\");\n\n\tglobal.Go = class {\n\t\tconstructor() {\n\t\t\tthis.argv = [\"js\"];\n\t\t\tthis.env = {};\n\t\t\tthis.exit = (code) => {\n\t\t\t\tif (code !== 0) {\n\t\t\t\t\tconsole.warn(\"exit code:\", code);\n\t\t\t\t}\n\t\t\t};\n\t\t\tthis._exitPromise = new Promise((resolve) => {\n\t\t\t\tthis._resolveExitPromise = resolve;\n\t\t\t});\n\t\t\tthis._pendingEvent = null;\n\t\t\tthis._scheduledTimeouts = new Map();\n\t\t\tthis._nextCallbackTimeoutID = 1;\n\n\t\t\tconst mem = () => {\n\t\t\t\t// The buffer may change when requesting more memory.\n\t\t\t\treturn new DataView(this._inst.exports.mem.buffer);\n\t\t\t}\n\n\t\t\tconst setInt64 = (addr, v) => {\n\t\t\t\tmem().setUint32(addr + 0, v, true);\n\t\t\t\tmem().setUint32(addr + 4, Math.floor(v / 4294967296), true);\n\t\t\t}\n\n\t\t\tconst getInt64 = (addr) => {\n\t\t\t\tconst low = mem().getUint32(addr + 0, true);\n\t\t\t\tconst high = mem().getInt32(addr + 4, true);\n\t\t\t\treturn low + high * 4294967296;\n\t\t\t}\n\n\t\t\tconst loadValue = (addr) => {\n\t\t\t\tconst f = mem().getFloat64(addr, true);\n\t\t\t\tif (f === 0) {\n\t\t\t\t\treturn undefined;\n\t\t\t\t}\n\t\t\t\tif (!isNaN(f)) {\n\t\t\t\t\treturn f;\n\t\t\t\t}\n\n\t\t\t\tconst id = mem().getUint32(addr, true);\n\t\t\t\treturn this._values[id];\n\t\t\t}\n\n\t\t\tconst storeValue = (addr, v) => {\n\t\t\t\tconst nanHead = 0x7FF80000;\n\n\t\t\t\tif (typeof v === \"number\") {\n\t\t\t\t\tif (isNaN(v)) {\n\t\t\t\t\t\tmem().setUint32(addr + 4, nanHead, true);\n\t\t\t\t\t\tmem().setUint32(addr, 0, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t\t}\n\t\t\t\t\tif (v === 0) {\n\t\t\t\t\t\tmem().setUint32(addr + 4, nanHead, true);\n\t\t\t\t\t\tmem().setUint32(addr, 1, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t\t}\n\t\t\t\t\tmem().setFloat64(addr, v, true);\n\t\t\t\t\treturn;\n\t\t\t\t}\n\n\t\t\t\tswitch (v) {\n\t\t\t\t\tcase undefined:\n\t\t\t\t\t\tmem().setFloat64(addr, 0, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t\tcase null:\n\t\t\t\t\t\tmem().setUint32(addr + 4, nanHead, true);\n\t\t\t\t\t\tmem().setUint32(addr, 2, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t\tcase true:\n\t\t\t\t\t\tmem().setUint32(addr + 4, nanHead, true);\n\t\t\t\t\t\tmem().setUint32(addr, 3, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t\tcase false:\n\t\t\t\t\t\tmem().setUint32(addr + 4, nanHead, true);\n\t\t\t\t\t\tmem().setUint32(addr, 4, true);\n\t\t\t\t\t\treturn;\n\t\t\t\t}\n\n\t\t\t\tlet ref = this._refs.get(v);\n\t\t\t\tif (ref === undefined) {\n\t\t\t\t\tref = this._values.length;\n\t\t\t\t\tthis._values.push(v);\n\t\t\t\t\tthis._refs.set(v, ref);\n\t\t\t\t}\n\t\t\t\tlet typeFlag = 0;\n\t\t\t\tswitch (typeof v) {\n\t\t\t\t\tcase \"string\":\n\t\t\t\t\t\ttypeFlag = 1;\n\t\t\t\t\t\tbreak;\n\t\t\t\t\tcase \"symbol\":\n\t\t\t\t\t\ttypeFlag = 2;\n\t\t\t\t\t\tbreak;\n\t\t\t\t\tcase \"function\":\n\t\t\t\t\t\ttypeFlag = 3;\n\t\t\t\t\t\tbreak;\n\t\t\t\t}\n\t\t\t\tmem().setUint32(addr + 4, nanHead | typeFlag, true);\n\t\t\t\tmem().setUint32(addr, ref, true);\n\t\t\t}\n\n\t\t\tconst loadSlice = (addr) => {\n\t\t\t\tconst array = getInt64(addr + 0);\n\t\t\t\tconst len = getInt64(addr + 8);\n\t\t\t\treturn new Uint8Array(this._inst.exports.mem.buffer, array, len);\n\t\t\t}\n\n\t\t\tconst loadSliceOfValues = (addr) => {\n\t\t\t\tconst array = getInt64(addr + 0);\n\t\t\t\tconst len = getInt64(addr + 8);\n\t\t\t\tconst a = new Array(len);\n\t\t\t\tfor (let i = 0; i < len; i++) {\n\t\t\t\t\ta[i] = loadValue(array + i * 8);\n\t\t\t\t}\n\t\t\t\treturn a;\n\t\t\t}\n\n\t\t\tconst loadString = (addr) => {\n\t\t\t\tconst saddr = getInt64(addr + 0);\n\t\t\t\tconst len = getInt64(addr + 8);\n\t\t\t\treturn decoder.decode(new DataView(this._inst.exports.mem.buffer, saddr, len));\n\t\t\t}\n\n\t\t\tconst timeOrigin = Date.now() - performance.now();\n\t\t\tthis.importObject = {\n\t\t\t\tgo: {\n\t\t\t\t\t// Go's SP does not change as long as no Go code is running. Some operations (e.g. calls, getters and setters)\n\t\t\t\t\t// may synchronously trigger a Go event handler. This makes Go code get executed in the middle of the imported\n\t\t\t\t\t// function. A goroutine can switch to a new stack if the current stack is too small (see morestack function).\n\t\t\t\t\t// This changes the SP, thus we have to update the SP used by the imported function.\n\n\t\t\t\t\t// func wasmExit(code int32)\n\t\t\t\t\t\"runtime.wasmExit\": (sp) => {\n\t\t\t\t\t\tconst code = mem().getInt32(sp + 8, true);\n\t\t\t\t\t\tthis.exited = true;\n\t\t\t\t\t\tdelete this._inst;\n\t\t\t\t\t\tdelete this._values;\n\t\t\t\t\t\tdelete this._refs;\n\t\t\t\t\t\tthis.exit(code);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)\n\t\t\t\t\t\"runtime.wasmWrite\": (sp) => {\n\t\t\t\t\t\tconst fd = getInt64(sp + 8);\n\t\t\t\t\t\tconst p = getInt64(sp + 16);\n\t\t\t\t\t\tconst n = mem().getInt32(sp + 24, true);\n\t\t\t\t\t\tfs.writeSync(fd, new Uint8Array(this._inst.exports.mem.buffer, p, n));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func nanotime() int64\n\t\t\t\t\t\"runtime.nanotime\": (sp) => {\n\t\t\t\t\t\tsetInt64(sp + 8, (timeOrigin + performance.now()) * 1000000);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func walltime() (sec int64, nsec int32)\n\t\t\t\t\t\"runtime.walltime\": (sp) => {\n\t\t\t\t\t\tconst msec = (new Date).getTime();\n\t\t\t\t\t\tsetInt64(sp + 8, msec / 1000);\n\t\t\t\t\t\tmem().setInt32(sp + 16, (msec % 1000) * 1000000, true);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func scheduleTimeoutEvent(delay int64) int32\n\t\t\t\t\t\"runtime.scheduleTimeoutEvent\": (sp) => {\n\t\t\t\t\t\tconst id = this._nextCallbackTimeoutID;\n\t\t\t\t\t\tthis._nextCallbackTimeoutID++;\n\t\t\t\t\t\tthis._scheduledTimeouts.set(id, setTimeout(\n\t\t\t\t\t\t\t() => {\n\t\t\t\t\t\t\t\tthis._resume();\n\t\t\t\t\t\t\t\twhile (this._scheduledTimeouts.has(id)) {\n\t\t\t\t\t\t\t\t\t// for some reason Go failed to register the timeout event, log and try again\n\t\t\t\t\t\t\t\t\t// (temporary workaround for https://github.com/golang/go/issues/28975)\n\t\t\t\t\t\t\t\t\tconsole.warn(\"scheduleTimeoutEvent: missed timeout event\");\n\t\t\t\t\t\t\t\t\tthis._resume();\n\t\t\t\t\t\t\t\t}\n\t\t\t\t\t\t\t},\n\t\t\t\t\t\t\tgetInt64(sp + 8) + 1, // setTimeout has been seen to fire up to 1 millisecond early\n\t\t\t\t\t\t));\n\t\t\t\t\t\tmem().setInt32(sp + 16, id, true);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func clearTimeoutEvent(id int32)\n\t\t\t\t\t\"runtime.clearTimeoutEvent\": (sp) => {\n\t\t\t\t\t\tconst id = mem().getInt32(sp + 8, true);\n\t\t\t\t\t\tclearTimeout(this._scheduledTimeouts.get(id));\n\t\t\t\t\t\tthis._scheduledTimeouts.delete(id);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func getRandomData(r []byte)\n\t\t\t\t\t\"runtime.getRandomData\": (sp) => {\n\t\t\t\t\t\tcrypto.getRandomValues(loadSlice(sp + 8));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func stringVal(value string) ref\n\t\t\t\t\t\"syscall/js.stringVal\": (sp) => {\n\t\t\t\t\t\tstoreValue(sp + 24, loadString(sp + 8));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueGet(v ref, p string) ref\n\t\t\t\t\t\"syscall/js.valueGet\": (sp) => {\n\t\t\t\t\t\tconst result = Reflect.get(loadValue(sp + 8), loadString(sp + 16));\n\t\t\t\t\t\tsp = this._inst.exports.getsp(); // see comment above\n\t\t\t\t\t\tstoreValue(sp + 32, result);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueSet(v ref, p string, x ref)\n\t\t\t\t\t\"syscall/js.valueSet\": (sp) => {\n\t\t\t\t\t\tReflect.set(loadValue(sp + 8), loadString(sp + 16), loadValue(sp + 32));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueIndex(v ref, i int) ref\n\t\t\t\t\t\"syscall/js.valueIndex\": (sp) => {\n\t\t\t\t\t\tstoreValue(sp + 24, Reflect.get(loadValue(sp + 8), getInt64(sp + 16)));\n\t\t\t\t\t},\n\n\t\t\t\t\t// valueSetIndex(v ref, i int, x ref)\n\t\t\t\t\t\"syscall/js.valueSetIndex\": (sp) => {\n\t\t\t\t\t\tReflect.set(loadValue(sp + 8), getInt64(sp + 16), loadValue(sp + 24));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueCall(v ref, m string, args []ref) (ref, bool)\n\t\t\t\t\t\"syscall/js.valueCall\": (sp) => {\n\t\t\t\t\t\ttry {\n\t\t\t\t\t\t\tconst v = loadValue(sp + 8);\n\t\t\t\t\t\t\tconst m = Reflect.get(v, loadString(sp + 16));\n\t\t\t\t\t\t\tconst args = loadSliceOfValues(sp + 32);\n\t\t\t\t\t\t\tconst result = Reflect.apply(m, v, args);\n\t\t\t\t\t\t\tsp = this._inst.exports.getsp(); // see comment above\n\t\t\t\t\t\t\tstoreValue(sp + 56, result);\n\t\t\t\t\t\t\tmem().setUint8(sp + 64, 1);\n\t\t\t\t\t\t} catch (err) {\n\t\t\t\t\t\t\tstoreValue(sp + 56, err);\n\t\t\t\t\t\t\tmem().setUint8(sp + 64, 0);\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueInvoke(v ref, args []ref) (ref, bool)\n\t\t\t\t\t\"syscall/js.valueInvoke\": (sp) => {\n\t\t\t\t\t\ttry {\n\t\t\t\t\t\t\tconst v = loadValue(sp + 8);\n\t\t\t\t\t\t\tconst args = loadSliceOfValues(sp + 16);\n\t\t\t\t\t\t\tconst result = Reflect.apply(v, undefined, args);\n\t\t\t\t\t\t\tsp = this._inst.exports.getsp(); // see comment above\n\t\t\t\t\t\t\tstoreValue(sp + 40, result);\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 1);\n\t\t\t\t\t\t} catch (err) {\n\t\t\t\t\t\t\tstoreValue(sp + 40, err);\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 0);\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueNew(v ref, args []ref) (ref, bool)\n\t\t\t\t\t\"syscall/js.valueNew\": (sp) => {\n\t\t\t\t\t\ttry {\n\t\t\t\t\t\t\tconst v = loadValue(sp + 8);\n\t\t\t\t\t\t\tconst args = loadSliceOfValues(sp + 16);\n\t\t\t\t\t\t\tconst result = Reflect.construct(v, args);\n\t\t\t\t\t\t\tsp = this._inst.exports.getsp(); // see comment above\n\t\t\t\t\t\t\tstoreValue(sp + 40, result);\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 1);\n\t\t\t\t\t\t} catch (err) {\n\t\t\t\t\t\t\tstoreValue(sp + 40, err);\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 0);\n\t\t\t\t\t\t}\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueLength(v ref) int\n\t\t\t\t\t\"syscall/js.valueLength\": (sp) => {\n\t\t\t\t\t\tsetInt64(sp + 16, parseInt(loadValue(sp + 8).length));\n\t\t\t\t\t},\n\n\t\t\t\t\t// valuePrepareString(v ref) (ref, int)\n\t\t\t\t\t\"syscall/js.valuePrepareString\": (sp) => {\n\t\t\t\t\t\tconst str = encoder.encode(String(loadValue(sp + 8)));\n\t\t\t\t\t\tstoreValue(sp + 16, str);\n\t\t\t\t\t\tsetInt64(sp + 24, str.length);\n\t\t\t\t\t},\n\n\t\t\t\t\t// valueLoadString(v ref, b []byte)\n\t\t\t\t\t\"syscall/js.valueLoadString\": (sp) => {\n\t\t\t\t\t\tconst str = loadValue(sp + 8);\n\t\t\t\t\t\tloadSlice(sp + 16).set(str);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func valueInstanceOf(v ref, t ref) bool\n\t\t\t\t\t\"syscall/js.valueInstanceOf\": (sp) => {\n\t\t\t\t\t\tmem().setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16));\n\t\t\t\t\t},\n\n\t\t\t\t\t// func copyBytesToGo(dst []byte, src ref) (int, bool)\n\t\t\t\t\t\"syscall/js.copyBytesToGo\": (sp) => {\n\t\t\t\t\t\tconst dst = loadSlice(sp + 8);\n\t\t\t\t\t\tconst src = loadValue(sp + 32);\n\t\t\t\t\t\tif (!(src instanceof Uint8Array)) {\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 0);\n\t\t\t\t\t\t\treturn;\n\t\t\t\t\t\t}\n\t\t\t\t\t\tconst toCopy = src.subarray(0, dst.length);\n\t\t\t\t\t\tdst.set(toCopy);\n\t\t\t\t\t\tsetInt64(sp + 40, toCopy.length);\n\t\t\t\t\t\tmem().setUint8(sp + 48, 1);\n\t\t\t\t\t},\n\n\t\t\t\t\t// func copyBytesToJS(dst ref, src []byte) (int, bool)\n\t\t\t\t\t\"syscall/js.copyBytesToJS\": (sp) => {\n\t\t\t\t\t\tconst dst = loadValue(sp + 8);\n\t\t\t\t\t\tconst src = loadSlice(sp + 16);\n\t\t\t\t\t\tif (!(dst instanceof Uint8Array)) {\n\t\t\t\t\t\t\tmem().setUint8(sp + 48, 0);\n\t\t\t\t\t\t\treturn;\n\t\t\t\t\t\t}\n\t\t\t\t\t\tconst toCopy = src.subarray(0, dst.length);\n\t\t\t\t\t\tdst.set(toCopy);\n\t\t\t\t\t\tsetInt64(sp + 40, toCopy.length);\n\t\t\t\t\t\tmem().setUint8(sp + 48, 1);\n\t\t\t\t\t},\n\n\t\t\t\t\t\"debug\": (value) => {\n\t\t\t\t\t\tconsole.log(value);\n\t\t\t\t\t},\n\t\t\t\t}\n\t\t\t};\n\t\t}\n\n\t\tasync run(instance) {\n\t\t\tthis._inst = instance;\n\t\t\tthis._values = [ // TODO: garbage collection\n\t\t\t\tNaN,\n\t\t\t\t0,\n\t\t\t\tnull,\n\t\t\t\ttrue,\n\t\t\t\tfalse,\n\t\t\t\tglobal,\n\t\t\t\tthis,\n\t\t\t];\n\t\t\tthis._refs = new Map();\n\t\t\tthis.exited = false;\n\n\t\t\tconst mem = new DataView(this._inst.exports.mem.buffer)\n\n\t\t\t// Pass command line arguments and environment variables to WebAssembly by writing them to the linear memory.\n\t\t\tlet offset = 4096;\n\n\t\t\tconst strPtr = (str) => {\n\t\t\t\tconst ptr = offset;\n\t\t\t\tconst bytes = encoder.encode(str + \"\\0\");\n\t\t\t\tnew Uint8Array(mem.buffer, offset, bytes.length).set(bytes);\n\t\t\t\toffset += bytes.length;\n\t\t\t\tif (offset % 8 !== 0) {\n\t\t\t\t\toffset += 8 - (offset % 8);\n\t\t\t\t}\n\t\t\t\treturn ptr;\n\t\t\t};\n\n\t\t\tconst argc = this.argv.length;\n\n\t\t\tconst argvPtrs = [];\n\t\t\tthis.argv.forEach((arg) => {\n\t\t\t\targvPtrs.push(strPtr(arg));\n\t\t\t});\n\n\t\t\tconst keys = Object.keys(this.env).sort();\n\t\t\targvPtrs.push(keys.length);\n\t\t\tkeys.forEach((key) => {\n\t\t\t\targvPtrs.push(strPtr(`${key}=${this.env[key]}`));\n\t\t\t});\n\n\t\t\tconst argv = offset;\n\t\t\targvPtrs.forEach((ptr) => {\n\t\t\t\tmem.setUint32(offset, ptr, true);\n\t\t\t\tmem.setUint32(offset + 4, 0, true);\n\t\t\t\toffset += 8;\n\t\t\t});\n\n\t\t\tthis._inst.exports.run(argc, argv);\n\t\t\tif (this.exited) {\n\t\t\t\tthis._resolveExitPromise();\n\t\t\t}\n\t\t\tawait this._exitPromise;\n\t\t}\n\n\t\t_resume() {\n\t\t\tif (this.exited) {\n\t\t\t\tthrow new Error(\"Go program has already exited\");\n\t\t\t}\n\t\t\tthis._inst.exports.resume();\n\t\t\tif (this.exited) {\n\t\t\t\tthis._resolveExitPromise();\n\t\t\t}\n\t\t}\n\n\t\t_makeFuncWrapper(id) {\n\t\t\tconst go = this;\n\t\t\treturn function () {\n\t\t\t\tconst event = { id: id, this: this, args: arguments };\n\t\t\t\tgo._pendingEvent = event;\n\t\t\t\tgo._resume();\n\t\t\t\treturn event.result;\n\t\t\t};\n\t\t}\n\t}\n\n\tif (\n\t\tglobal.require &&\n\t\tglobal.require.main === module &&\n\t\tglobal.process &&\n\t\tglobal.process.versions &&\n\t\t!global.process.versions.electron\n\t) {\n\t\tif (process.argv.length < 3) {\n\t\t\tconsole.error(\"usage: go_js_wasm_exec [wasm binary] [arguments]\");\n\t\t\tprocess.exit(1);\n\t\t}\n\n\t\tconst go = new Go();\n\t\tgo.argv = process.argv.slice(2);\n\t\tgo.env = Object.assign({ TMPDIR: require(\"os\").tmpdir() }, process.env);\n\t\tgo.exit = process.exit;\n\t\tWebAssembly.instantiate(fs.readFileSync(process.argv[2]), go.importObject).then((result) => {\n\t\t\tprocess.on(\"exit\", (code) => { // Node.js exits if no event handler is pending\n\t\t\t\tif (code === 0 && !go.exited) {\n\t\t\t\t\t// deadlock, make Go print error and stack traces\n\t\t\t\t\tgo._pendingEvent = { id: 0 };\n\t\t\t\t\tgo._resume();\n\t\t\t\t}\n\t\t\t});\n\t\t\treturn go.run(result.instance);\n\t\t}).catch((err) => {\n\t\t\tconsole.error(err);\n\t\t\tprocess.exit(1);\n\t\t});\n\t}\n})();\n"

	viJS = "// -----------------------------------------------------------------------------\n// Init service worker\n// -----------------------------------------------------------------------------\nif (\"serviceWorker\" in navigator) {\n  navigator.serviceWorker\n    .register(\"/app-worker.js\")\n    .then(reg => {\n      console.log(\"registering app service worker\");\n    })\n    .catch(err => {\n      console.error(\"offline service worker registration failed\", err);\n    });\n}\n\n// -----------------------------------------------------------------------------\n// Init progressive app\n// -----------------------------------------------------------------------------\nlet deferredPrompt;\n\nwindow.addEventListener(\"beforeinstallprompt\", e => {\n  e.preventDefault();\n  deferredPrompt = e;\n});\n\n// -----------------------------------------------------------------------------\n// Init Web Assembly\n// -----------------------------------------------------------------------------\nif (!WebAssembly.instantiateStreaming) {\n  WebAssembly.instantiateStreaming = async (resp, importObject) => {\n    const source = await (await resp).arrayBuffer();\n    return await WebAssembly.instantiate(source, importObject);\n  };\n}\n\nconst go = new Go();\n\nWebAssembly.instantiateStreaming(fetch(\"{{.Wasm}}\"), go.importObject)\n  .then(result => {\n    go.run(result.instance);\n  })\n  .catch(err => {\n    const loaderIcon = document.getElementById(\"app-wasm-loader-icon\");\n    loaderIcon.className = \"\";\n\n    const loaderLabel = document.getElementById(\"app-wasm-loader-label\");\n    loadingLabel.innerText = err;\n\n    console.error(\"loading wasm failed: \" + err);\n  });\n"

	viWorkerJS = "const cacheName = \"app-\" + \"{{.Version}}\";\n\nself.addEventListener(\"install\", event => {\n  console.log(\"installing app worker\");\n  self.skipWaiting();\n\n  event.waitUntil(\n    caches.open(cacheName).then(cache => {\n      return cache.addAll([\n        {{range $path, $element := .ResourcesToCache}}\"{{$path}}\",\n        {{end}}\n      ]);\n    })\n  );\n});\n\nself.addEventListener(\"activate\", event => {\n  console.log(\"app worker is activating\");\n  event.waitUntil(\n    caches.keys().then(keyList => {\n      return Promise.all(\n        keyList.map(key => {\n          if (key !== cacheName) {\n            return caches.delete(key);\n          }\n        })\n      );\n    })\n  );\n});\n\nself.addEventListener(\"fetch\", event => {\n  event.respondWith(\n    caches.match(event.request).then(response => {\n      return response || fetch(event.request);\n    })\n  );\n});\n"

	manifestJSON = "{\n  \"short_name\": \"{{.ShortName}}\",\n  \"name\": \"{{.Name}}\",\n  \"icons\": [\n    {\n      \"src\": \"{{.DefaultIcon}}\",\n      \"type\": \"image/png\",\n      \"sizes\": \"192x192\"\n    },\n    {\n      \"src\": \"{{.LargeIcon}}\",\n      \"type\": \"image/png\",\n      \"sizes\": \"512x512\"\n    }\n  ],\n  \"start_url\": \"/\",\n  \"background_color\": \"{{.BackgroundColor}}\",\n  \"display\": \"standalone\",\n  \"theme_color\": \"{{.ThemeColor}}\"\n}\n"

	viCSS = "/* Basics */\nhtml {\n  height: 100%;\n  width: 100%;\n  margin: 0;\n  padding: 0;\n  overflow: hidden;\n}\n\nbody {\n  height: 100%;\n  width: 100%;\n  margin: 0;\n  padding: 0;\n  overflow: hidden;\n  font-family: -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, Oxygen,\n    Ubuntu, Cantarell, \"Open Sans\", \"Helvetica Neue\", sans-serif;\n  font-size: 10pt;\n  font-weight: 400;\n  color: white;\n  background-color: #2d2c2c;\n}\n\n@media (prefers-color-scheme: light) {\n  body {\n    color: black;\n    background-color: #f6f6f6;\n  }\n}\n\nh1 {\n  width: calc(100% - 24px - 24px);\n  margin: 0;\n  padding: 0 24px 12px;\n  outline: 0;\n  font-size: 2.7em;\n  font-weight: 200;\n  letter-spacing: 1px;\n}\n\nh2,\nh3,\nh4,\nh5,\nh6 {\n  width: calc(100% - 24px - 24px);\n  margin: 0;\n  padding: 0 24px 12px;\n  outline: 0;\n  font-weight: 200;\n}\n\np {\n  width: calc(100% - 24px - 24px);\n  margin: 0;\n  padding: 0 24px 12px;\n  outline: 0;\n}\n\na {\n  color: currentColor;\n  text-decoration: none;\n  cursor: pointer;\n}\n\na:hover {\n  color: deepskyblue;\n}\n\ninput {\n  width: calc(100% - 24px - 24px);\n  margin: 0 24px 12px;\n  padding: 0 0 3px;\n  border: 0;\n  border-radius: 0;\n  border-bottom: 1px solid currentColor;\n  outline: none;\n  box-shadow: none;\n  display: block;\n  background-color: transparent;\n  color: currentColor;\n  font-size: inherit;\n  font-family: inherit;\n}\n\ninput:hover,\ninput:focus {\n  border-bottom-color: deepskyblue;\n  caret-color: deepskyblue;\n}\n\ninput[type=\"number\"]::-webkit-inner-spin-button,\ninput[type=\"number\"]::-webkit-outer-spin-button {\n  -webkit-appearance: none;\n  margin: 0;\n}\n\nul {\n  margin: 0;\n  padding: 3px 24px 6px 42px;\n  outline: 0;\n}\n\nul li {\n  margin: 0;\n  padding: 3px 0;\n}\n\nul li:first-child {\n  padding: 0 0 3px;\n}\n\nul li:last-child {\n  padding: 3px 0 0;\n}\n\ntable {\n  width: calc(100% - 24px - 24px);\n  margin: 0 24px;\n  border-collapse: collapse;\n  table-layout: fixed;\n  font-size: inherit;\n}\n\ntable th {\n  padding: 0 12px 12px 0;\n  border-bottom: 1px solid currentColor;\n  font-size: 12pt;\n  font-weight: 200;\n  text-align: start;\n}\n\ntable td {\n  padding: 12px 12px 12px 0;\n  border-bottom: 1px solid currentColor;\n}\n\ntable tr:first-child td {\n  padding: 0 12px 12px 0;\n}\n\ntable tr:last-child td {\n  border-bottom: 0;\n  padding: 12px 12px 12px 0;\n}\n\nbutton,\n.app-button {\n  min-width: 100px;\n  margin: 0 6px 12px 0;\n  padding: 8px;\n  border: 1px solid currentColor;\n  border-radius: 8px;\n  outline: none;\n  display: inline-block;\n  background: none;\n  color: inherit;\n  font: inherit;\n  font-size: inherit;\n  text-align: center;\n  cursor: pointer;\n  -webkit-touch-callout: none;\n  -webkit-user-select: none;\n  -khtml-user-select: none;\n  -moz-user-select: none;\n  -ms-user-select: none;\n}\n\n::-webkit-scrollbar {\n  background-color: rgba(255, 255, 255, 0.04);\n}\n\n::-webkit-scrollbar-thumb {\n  background-color: rgba(255, 255, 255, 0.05);\n}\n\n/* Loader and Not Found */\n.app-wasm-layout {\n  display: flex;\n  flex-direction: column;\n  justify-content: center;\n  align-items: center;\n  width: 100%;\n  height: 100%;\n}\n\n.app-wasm-icon {\n  max-width: 100px;\n  max-height: 100px;\n  user-select: none;\n  margin: 0 6px;\n  -moz-user-select: none;\n  -webkit-user-drag: none;\n  -webkit-user-select: none;\n  -ms-user-select: none;\n}\n\n.app-wasm-label {\n  margin-top: 6px;\n  font-size: 16pt;\n  font-weight: 100;\n  letter-spacing: 1px;\n  max-width: 480px;\n  text-align: center;\n  text-transform: lowercase;\n}\n\n.app-notfound-title {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n  font-size: 65pt;\n  font-weight: 100;\n}\n\n.app-spin {\n  animation: app-spin-frames 2.5s infinite linear;\n}\n\n@keyframes app-spin-frames {\n  from {\n    transform: rotate(0deg);\n  }\n\n  to {\n    transform: rotate(360deg);\n  }\n}\n\n/* Menu */\n#app-contextmenu-background {\n  position: fixed;\n  width: 100%;\n  height: 100%;\n  overflow: hidden;\n  left: 0;\n  top: 0;\n  z-index: 42000;\n}\n\n.app-contextmenu {\n  position: fixed;\n  min-width: 150px;\n  max-width: 480px;\n  padding: 6px 0;\n  border-radius: 6px;\n  border: solid 1px rgba(255, 255, 255, 0.1);\n  background-color: rgba(40, 38, 37, 0.97);\n  color: white;\n  -webkit-box-shadow: -1px 12px 38px 0px rgba(0, 0, 0, 0.6);\n  -moz-box-shadow: -1px 12px 38px 0px rgba(0, 0, 0, 0.6);\n  box-shadow: -1px 12px 38px 0px rgba(0, 0, 0, 0.6);\n}\n\n.app-contextmenu-hidden {\n  display: none;\n}\n\n.app-contextmenu-visible {\n  display: block;\n}\n\n@media (prefers-color-scheme: light) {\n  .app-contextmenu {\n    color: black;\n    background-color: rgba(221, 221, 221, 0.97);\n    border: solid 1px rgba(0, 0, 0, 0.2);\n  }\n}\n\n.app-menuitem {\n  display: flex;\n  align-items: center;\n  border: 0;\n  border-radius: 0;\n  margin: 0;\n  padding: 3px 24px;\n  width: 100%;\n  text-align: left;\n}\n\n.app-menuitem:hover {\n  background-color: deepskyblue;\n}\n\n.app-menuitem:disabled {\n  opacity: 0.15;\n  background-color: transparent;\n}\n\n.app-menuitem-separator {\n  width: 100%;\n  height: 0;\n  margin: 6px 0;\n  border-top: solid 1px rgba(255, 255, 255, 0.1);\n}\n\n@media (prefers-color-scheme: light) {\n  .app-menuitem-separator {\n    border-top: solid 1px rgba(0, 0, 0, 0.2);\n  }\n}\n\n.app-menuitem-label {\n  user-select: none;\n  flex-grow: 1;\n}\n\n.app-menuitem-keys {\n  flex-grow: 0;\n  margin-left: 12px;\n  text-transform: capitalize;\n}\n\n.app-menuitem-icon {\n  width: 18px;\n  height: 18px;\n  margin-right: 12px;\n}\n"
)
