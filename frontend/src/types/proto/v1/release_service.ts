/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "bytebase.v1";

export interface GetReleaseRequest {
  /** Format: projects/{project}/releases/{release} */
  name: string;
}

export interface ListReleasesRequest {
  /** Format: projects/{project} */
  parent: string;
  /**
   * The maximum number of change histories to return. The service may return fewer than this value.
   * If unspecified, at most 10 change histories will be returned.
   * The maximum value is 1000; values above 1000 will be coerced to 1000.
   */
  pageSize: number;
  /**
   * Not used. A page token, received from a previous `ListReleasesRequest` call.
   * Provide this to retrieve the subsequent page.
   *
   * When paginating, all other parameters provided to `ListReleasesRequest` must match
   * the call that provided the page token.
   */
  pageToken: string;
}

export interface ListReleasesResponse {
  releases: Release[];
  /**
   * A token, which can be sent as `page_token` to retrieve the next page.
   * If this field is omitted, there are no subsequent pages.
   */
  nextPageToken: string;
}

export interface Release {
  /** Format: projects/{project}/releases/{release} */
  name: string;
  files: File[];
  vcsSource: Release_VCSSource | undefined;
}

export interface Release_VCSSource {
  pullRequestUrl: string;
}

export interface File {
  filename: string;
  /** sheet holds the statement. */
  sheet: string;
  sheetHash: Uint8Array;
  type: File_Type;
  version: string;
}

export enum File_Type {
  VERSIONED = "VERSIONED",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function file_TypeFromJSON(object: any): File_Type {
  switch (object) {
    case 0:
    case "VERSIONED":
      return File_Type.VERSIONED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return File_Type.UNRECOGNIZED;
  }
}

export function file_TypeToJSON(object: File_Type): string {
  switch (object) {
    case File_Type.VERSIONED:
      return "VERSIONED";
    case File_Type.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export function file_TypeToNumber(object: File_Type): number {
  switch (object) {
    case File_Type.VERSIONED:
      return 0;
    case File_Type.UNRECOGNIZED:
    default:
      return -1;
  }
}

function createBaseGetReleaseRequest(): GetReleaseRequest {
  return { name: "" };
}

export const GetReleaseRequest = {
  encode(message: GetReleaseRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetReleaseRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetReleaseRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetReleaseRequest {
    return { name: isSet(object.name) ? globalThis.String(object.name) : "" };
  },

  toJSON(message: GetReleaseRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create(base?: DeepPartial<GetReleaseRequest>): GetReleaseRequest {
    return GetReleaseRequest.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<GetReleaseRequest>): GetReleaseRequest {
    const message = createBaseGetReleaseRequest();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseListReleasesRequest(): ListReleasesRequest {
  return { parent: "", pageSize: 0, pageToken: "" };
}

export const ListReleasesRequest = {
  encode(message: ListReleasesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.parent !== "") {
      writer.uint32(10).string(message.parent);
    }
    if (message.pageSize !== 0) {
      writer.uint32(16).int32(message.pageSize);
    }
    if (message.pageToken !== "") {
      writer.uint32(26).string(message.pageToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListReleasesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListReleasesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.parent = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.pageSize = reader.int32();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.pageToken = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListReleasesRequest {
    return {
      parent: isSet(object.parent) ? globalThis.String(object.parent) : "",
      pageSize: isSet(object.pageSize) ? globalThis.Number(object.pageSize) : 0,
      pageToken: isSet(object.pageToken) ? globalThis.String(object.pageToken) : "",
    };
  },

  toJSON(message: ListReleasesRequest): unknown {
    const obj: any = {};
    if (message.parent !== "") {
      obj.parent = message.parent;
    }
    if (message.pageSize !== 0) {
      obj.pageSize = Math.round(message.pageSize);
    }
    if (message.pageToken !== "") {
      obj.pageToken = message.pageToken;
    }
    return obj;
  },

  create(base?: DeepPartial<ListReleasesRequest>): ListReleasesRequest {
    return ListReleasesRequest.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<ListReleasesRequest>): ListReleasesRequest {
    const message = createBaseListReleasesRequest();
    message.parent = object.parent ?? "";
    message.pageSize = object.pageSize ?? 0;
    message.pageToken = object.pageToken ?? "";
    return message;
  },
};

function createBaseListReleasesResponse(): ListReleasesResponse {
  return { releases: [], nextPageToken: "" };
}

export const ListReleasesResponse = {
  encode(message: ListReleasesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.releases) {
      Release.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.nextPageToken !== "") {
      writer.uint32(18).string(message.nextPageToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListReleasesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListReleasesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.releases.push(Release.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.nextPageToken = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListReleasesResponse {
    return {
      releases: globalThis.Array.isArray(object?.releases) ? object.releases.map((e: any) => Release.fromJSON(e)) : [],
      nextPageToken: isSet(object.nextPageToken) ? globalThis.String(object.nextPageToken) : "",
    };
  },

  toJSON(message: ListReleasesResponse): unknown {
    const obj: any = {};
    if (message.releases?.length) {
      obj.releases = message.releases.map((e) => Release.toJSON(e));
    }
    if (message.nextPageToken !== "") {
      obj.nextPageToken = message.nextPageToken;
    }
    return obj;
  },

  create(base?: DeepPartial<ListReleasesResponse>): ListReleasesResponse {
    return ListReleasesResponse.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<ListReleasesResponse>): ListReleasesResponse {
    const message = createBaseListReleasesResponse();
    message.releases = object.releases?.map((e) => Release.fromPartial(e)) || [];
    message.nextPageToken = object.nextPageToken ?? "";
    return message;
  },
};

function createBaseRelease(): Release {
  return { name: "", files: [], vcsSource: undefined };
}

export const Release = {
  encode(message: Release, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    for (const v of message.files) {
      File.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.vcsSource !== undefined) {
      Release_VCSSource.encode(message.vcsSource, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Release {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRelease();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.files.push(File.decode(reader, reader.uint32()));
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.vcsSource = Release_VCSSource.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Release {
    return {
      name: isSet(object.name) ? globalThis.String(object.name) : "",
      files: globalThis.Array.isArray(object?.files) ? object.files.map((e: any) => File.fromJSON(e)) : [],
      vcsSource: isSet(object.vcsSource) ? Release_VCSSource.fromJSON(object.vcsSource) : undefined,
    };
  },

  toJSON(message: Release): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.files?.length) {
      obj.files = message.files.map((e) => File.toJSON(e));
    }
    if (message.vcsSource !== undefined) {
      obj.vcsSource = Release_VCSSource.toJSON(message.vcsSource);
    }
    return obj;
  },

  create(base?: DeepPartial<Release>): Release {
    return Release.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<Release>): Release {
    const message = createBaseRelease();
    message.name = object.name ?? "";
    message.files = object.files?.map((e) => File.fromPartial(e)) || [];
    message.vcsSource = (object.vcsSource !== undefined && object.vcsSource !== null)
      ? Release_VCSSource.fromPartial(object.vcsSource)
      : undefined;
    return message;
  },
};

function createBaseRelease_VCSSource(): Release_VCSSource {
  return { pullRequestUrl: "" };
}

export const Release_VCSSource = {
  encode(message: Release_VCSSource, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pullRequestUrl !== "") {
      writer.uint32(10).string(message.pullRequestUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Release_VCSSource {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRelease_VCSSource();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.pullRequestUrl = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Release_VCSSource {
    return { pullRequestUrl: isSet(object.pullRequestUrl) ? globalThis.String(object.pullRequestUrl) : "" };
  },

  toJSON(message: Release_VCSSource): unknown {
    const obj: any = {};
    if (message.pullRequestUrl !== "") {
      obj.pullRequestUrl = message.pullRequestUrl;
    }
    return obj;
  },

  create(base?: DeepPartial<Release_VCSSource>): Release_VCSSource {
    return Release_VCSSource.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<Release_VCSSource>): Release_VCSSource {
    const message = createBaseRelease_VCSSource();
    message.pullRequestUrl = object.pullRequestUrl ?? "";
    return message;
  },
};

function createBaseFile(): File {
  return { filename: "", sheet: "", sheetHash: new Uint8Array(0), type: File_Type.VERSIONED, version: "" };
}

export const File = {
  encode(message: File, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.filename !== "") {
      writer.uint32(10).string(message.filename);
    }
    if (message.sheet !== "") {
      writer.uint32(18).string(message.sheet);
    }
    if (message.sheetHash.length !== 0) {
      writer.uint32(26).bytes(message.sheetHash);
    }
    if (message.type !== File_Type.VERSIONED) {
      writer.uint32(32).int32(file_TypeToNumber(message.type));
    }
    if (message.version !== "") {
      writer.uint32(42).string(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): File {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFile();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.filename = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.sheet = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.sheetHash = reader.bytes();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.type = file_TypeFromJSON(reader.int32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.version = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): File {
    return {
      filename: isSet(object.filename) ? globalThis.String(object.filename) : "",
      sheet: isSet(object.sheet) ? globalThis.String(object.sheet) : "",
      sheetHash: isSet(object.sheetHash) ? bytesFromBase64(object.sheetHash) : new Uint8Array(0),
      type: isSet(object.type) ? file_TypeFromJSON(object.type) : File_Type.VERSIONED,
      version: isSet(object.version) ? globalThis.String(object.version) : "",
    };
  },

  toJSON(message: File): unknown {
    const obj: any = {};
    if (message.filename !== "") {
      obj.filename = message.filename;
    }
    if (message.sheet !== "") {
      obj.sheet = message.sheet;
    }
    if (message.sheetHash.length !== 0) {
      obj.sheetHash = base64FromBytes(message.sheetHash);
    }
    if (message.type !== File_Type.VERSIONED) {
      obj.type = file_TypeToJSON(message.type);
    }
    if (message.version !== "") {
      obj.version = message.version;
    }
    return obj;
  },

  create(base?: DeepPartial<File>): File {
    return File.fromPartial(base ?? {});
  },
  fromPartial(object: DeepPartial<File>): File {
    const message = createBaseFile();
    message.filename = object.filename ?? "";
    message.sheet = object.sheet ?? "";
    message.sheetHash = object.sheetHash ?? new Uint8Array(0);
    message.type = object.type ?? File_Type.VERSIONED;
    message.version = object.version ?? "";
    return message;
  },
};

export type ReleaseServiceDefinition = typeof ReleaseServiceDefinition;
export const ReleaseServiceDefinition = {
  name: "ReleaseService",
  fullName: "bytebase.v1.ReleaseService",
  methods: {
    getRelease: {
      name: "GetRelease",
      requestType: GetReleaseRequest,
      requestStream: false,
      responseType: Release,
      responseStream: false,
      options: {
        _unknownFields: {
          8410: [new Uint8Array([4, 110, 97, 109, 101])],
          578365826: [
            new Uint8Array([
              34,
              18,
              32,
              47,
              118,
              49,
              47,
              123,
              110,
              97,
              109,
              101,
              61,
              112,
              114,
              111,
              106,
              101,
              99,
              116,
              115,
              47,
              42,
              47,
              114,
              101,
              108,
              101,
              97,
              115,
              101,
              115,
              47,
              42,
              125,
            ]),
          ],
        },
      },
    },
    listReleases: {
      name: "ListReleases",
      requestType: ListReleasesRequest,
      requestStream: false,
      responseType: ListReleasesResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8410: [new Uint8Array([6, 112, 97, 114, 101, 110, 116])],
          578365826: [
            new Uint8Array([
              34,
              18,
              32,
              47,
              118,
              49,
              47,
              123,
              112,
              97,
              114,
              101,
              110,
              116,
              61,
              112,
              114,
              111,
              106,
              101,
              99,
              116,
              115,
              47,
              42,
              125,
              47,
              114,
              101,
              108,
              101,
              97,
              115,
              101,
              115,
            ]),
          ],
        },
      },
    },
  },
} as const;

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(globalThis.String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Long ? string | number | Long : T extends globalThis.Array<infer U> ? globalThis.Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
