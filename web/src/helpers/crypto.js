/* eslint-disable */
const SHA256 = (function() {
    function u(c, a, b) {
      var g = 0,
        d = [],
        f = 0,
        e,
        h,
        q,
        m,
        n,
        x,
        t,
        p = !1,
        l = [],
        r = [],
        k,
        y = !1;
      b = b || {};
      e = b.encoding || "UTF8";
      k = b.numRounds || 1;
      q = z(a, e);
      if (k !== parseInt(k, 10) || 1 > k)
        throw Error("numRounds must a integer >= 1");
      x = function(a, b) {
        return A(a, b, c);
      };
      t = function(a, b, f, d) {
        var g, e;
        if ("SHA-224" === c || "SHA-256" === c)
          (g = (((b + 65) >>> 9) << 4) + 15), (e = 16);
        else throw Error("Unexpected error in SHA-2 implementation");
        for (; a.length <= g; ) a.push(0);
        a[b >>> 5] |= 128 << (24 - (b % 32));
        b = b + f;
        a[g] = b & 4294967295;
        a[g - 1] = (b / 4294967296) | 0;
        f = a.length;
        for (b = 0; b < f; b += e) d = A(a.slice(b, b + e), d, c);
        if ("SHA-224" === c) a = [d[0], d[1], d[2], d[3], d[4], d[5], d[6]];
        else if ("SHA-256" === c) a = d;
        else throw Error("Unexpected error in SHA-2 implementation");
        return a;
      };
      if ("SHA-224" === c) (n = 512), (m = 224);
      else if ("SHA-256" === c) (n = 512), (m = 256);
      else throw Error("Chosen SHA variant is not supported");
      h = v(c);
      this.setHMACKey = function(a, b, d) {
        var f;
        if (!0 === p) throw Error("HMAC key already set");
        if (!0 === y) throw Error("Cannot set HMAC key after calling update");
        e = (d || {}).encoding || "UTF8";
        b = z(b, e)(a);
        a = b.binLen;
        b = b.value;
        f = n >>> 3;
        d = f / 4 - 1;
        if (f < a / 8) {
          for (b = t(b, a, 0, v(c)); b.length <= d; ) b.push(0);
          b[d] &= 4294967040;
        } else if (f > a / 8) {
          for (; b.length <= d; ) b.push(0);
          b[d] &= 4294967040;
        }
        for (a = 0; a <= d; a += 1)
          (l[a] = b[a] ^ 909522486), (r[a] = b[a] ^ 1549556828);
        h = x(l, h);
        g = n;
        p = !0;
      };
      this.update = function(a) {
        var b,
          c,
          e,
          m = 0,
          p = n >>> 5;
        b = q(a, d, f);
        a = b.binLen;
        c = b.value;
        b = a >>> 5;
        for (e = 0; e < b; e += p)
          m + n <= a && ((h = x(c.slice(e, e + p), h)), (m += n));
        g += m;
        d = c.slice(m >>> 5);
        f = a % n;
        y = !0;
      };
      this.getHash = function(a, b) {
        var e, n, q, l;
        if (!0 === p) throw Error("Cannot call getHash after setting HMAC key");
        q = B(b);
        switch (a) {
          case "HEX":
            e = function(a) {
              return C(a, q);
            };
            break;
          case "B64":
            e = function(a) {
              return D(a, q);
            };
            break;
          case "BYTES":
            e = E;
            break;
          case "ARRAYBUFFER":
            try {
              n = new ArrayBuffer(0);
            } catch (x) {
              throw Error("ARRAYBUFFER not supported by this environment");
            }
            e = F;
            break;
          default:
            throw Error("format must be HEX, B64, BYTES, or ARRAYBUFFER");
        }
        l = t(d.slice(), f, g, h.slice());
        for (n = 1; n < k; n += 1) l = t(l, m, 0, v(c));
        return e(l);
      };
      this.getHMAC = function(a, b) {
        var e, q, l, k;
        if (!1 === p)
          throw Error("Cannot call getHMAC without first setting HMAC key");
        l = B(b);
        switch (a) {
          case "HEX":
            e = function(a) {
              return C(a, l);
            };
            break;
          case "B64":
            e = function(a) {
              return D(a, l);
            };
            break;
          case "BYTES":
            e = E;
            break;
          case "ARRAYBUFFER":
            try {
              e = new ArrayBuffer(0);
            } catch (G) {
              throw Error("ARRAYBUFFER not supported by this environment");
            }
            e = F;
            break;
          default:
            throw Error("outputFormat must be HEX, B64, BYTES, or ARRAYBUFFER");
        }
        q = t(d.slice(), f, g, h.slice());
        k = x(r, v(c));
        k = t(q, m, n, k);
        return e(k);
      };
    }
  
    function l() {}
  
    function J(c, a, b) {
      var g = c.length,
        d,
        f,
        e,
        h,
        q;
      a = a || [0];
      b = b || 0;
      q = b >>> 3;
      if (0 !== g % 2)
        throw Error("String of HEX type must be in byte increments");
      for (d = 0; d < g; d += 2) {
        f = parseInt(c.substr(d, 2), 16);
        if (isNaN(f))
          throw Error("String of HEX type contains invalid characters");
        h = (d >>> 1) + q;
        for (e = h >>> 2; a.length <= e; ) a.push(0);
        a[e] |= f << (8 * (3 - (h % 4)));
      }
      return {
        value: a,
        binLen: 4 * g + b
      };
    }
  
    function K(c, a, b) {
      var g = [],
        d,
        f,
        e,
        h,
        g = a || [0];
      b = b || 0;
      f = b >>> 3;
      for (d = 0; d < c.length; d += 1)
        (a = c.charCodeAt(d)),
          (h = d + f),
          (e = h >>> 2),
          g.length <= e && g.push(0),
          (g[e] |= a << (8 * (3 - (h % 4))));
      return {
        value: g,
        binLen: 8 * c.length + b
      };
    }
  
    function L(c, a, b) {
      var g = [],
        d = 0,
        f,
        e,
        h,
        q,
        m,
        n,
        g = a || [0];
      b = b || 0;
      a = b >>> 3;
      if (-1 === c.search(/^[a-zA-Z0-9=+\/]+$/))
        throw Error("Invalid character in base-64 string");
      e = c.indexOf("=");
      c = c.replace(/\=/g, "");
      if (-1 !== e && e < c.length)
        throw Error("Invalid '=' found in base-64 string");
      for (e = 0; e < c.length; e += 4) {
        m = c.substr(e, 4);
        for (h = q = 0; h < m.length; h += 1)
          (f = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".indexOf(
            m[h]
          )),
            (q |= f << (18 - 6 * h));
        for (h = 0; h < m.length - 1; h += 1) {
          n = d + a;
          for (f = n >>> 2; g.length <= f; ) g.push(0);
          g[f] |= ((q >>> (16 - 8 * h)) & 255) << (8 * (3 - (n % 4)));
          d += 1;
        }
      }
      return {
        value: g,
        binLen: 8 * d + b
      };
    }
  
    function M(c, a, b) {
      var g = [],
        d,
        f,
        e,
        g = a || [0];
      b = b || 0;
      d = b >>> 3;
      for (a = 0; a < c.byteLength; a += 1)
        (e = a + d),
          (f = e >>> 2),
          g.length <= f && g.push(0),
          (g[f] |= c[a] << (8 * (3 - (e % 4))));
      return {
        value: g,
        binLen: 8 * c.byteLength + b
      };
    }
  
    function C(c, a) {
      var b = "",
        g = 4 * c.length,
        d,
        f;
      for (d = 0; d < g; d += 1)
        (f = c[d >>> 2] >>> (8 * (3 - (d % 4)))),
          (b +=
            "0123456789abcdef".charAt((f >>> 4) & 15) +
            "0123456789abcdef".charAt(f & 15));
      return a.outputUpper ? b.toUpperCase() : b;
    }
  
    function D(c, a) {
      var b = "",
        g = 4 * c.length,
        d,
        f,
        e;
      for (d = 0; d < g; d += 3)
        for (
          e = (d + 1) >>> 2,
            f = c.length <= e ? 0 : c[e],
            e = (d + 2) >>> 2,
            e = c.length <= e ? 0 : c[e],
            e =
              (((c[d >>> 2] >>> (8 * (3 - (d % 4)))) & 255) << 16) |
              (((f >>> (8 * (3 - ((d + 1) % 4)))) & 255) << 8) |
              ((e >>> (8 * (3 - ((d + 2) % 4)))) & 255),
            f = 0;
          4 > f;
          f += 1
        )
          8 * d + 6 * f <= 32 * c.length
            ? (b += "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/".charAt(
                (e >>> (6 * (3 - f))) & 63
              ))
            : (b += a.b64Pad);
      return b;
    }
  
    function E(c) {
      var a = "",
        b = 4 * c.length,
        g,
        d;
      for (g = 0; g < b; g += 1)
        (d = (c[g >>> 2] >>> (8 * (3 - (g % 4)))) & 255),
          (a += String.fromCharCode(d));
      return a;
    }
  
    function F(c) {
      var a = 4 * c.length,
        b,
        g = new ArrayBuffer(a);
      for (b = 0; b < a; b += 1)
        g[b] = (c[b >>> 2] >>> (8 * (3 - (b % 4)))) & 255;
      return g;
    }
  
    function B(c) {
      var a = {
        outputUpper: !1,
        b64Pad: "="
      };
      c = c || {};
      a.outputUpper = c.outputUpper || !1;
      !0 === c.hasOwnProperty("b64Pad") && (a.b64Pad = c.b64Pad);
      if ("boolean" !== typeof a.outputUpper)
        throw Error("Invalid outputUpper formatting option");
      if ("string" !== typeof a.b64Pad)
        throw Error("Invalid b64Pad formatting option");
      return a;
    }
  
    function z(c, a) {
      var b;
      switch (a) {
        case "UTF8":
        case "UTF16BE":
        case "UTF16LE":
          break;
        default:
          throw Error("encoding must be UTF8, UTF16BE, or UTF16LE");
      }
      switch (c) {
        case "HEX":
          b = J;
          break;
        case "TEXT":
          b = function(b, c, f) {
            var e = [],
              h = [],
              q = 0,
              m,
              n,
              l,
              k,
              p,
              e = c || [0];
            c = f || 0;
            l = c >>> 3;
            if ("UTF8" === a)
              for (m = 0; m < b.length; m += 1)
                for (
                  f = b.charCodeAt(m),
                    h = [],
                    128 > f
                      ? h.push(f)
                      : 2048 > f
                      ? (h.push(192 | (f >>> 6)), h.push(128 | (f & 63)))
                      : 55296 > f || 57344 <= f
                      ? h.push(
                          224 | (f >>> 12),
                          128 | ((f >>> 6) & 63),
                          128 | (f & 63)
                        )
                      : ((m += 1),
                        (f =
                          65536 +
                          (((f & 1023) << 10) | (b.charCodeAt(m) & 1023))),
                        h.push(
                          240 | (f >>> 18),
                          128 | ((f >>> 12) & 63),
                          128 | ((f >>> 6) & 63),
                          128 | (f & 63)
                        )),
                    n = 0;
                  n < h.length;
                  n += 1
                ) {
                  p = q + l;
                  for (k = p >>> 2; e.length <= k; ) e.push(0);
                  e[k] |= h[n] << (8 * (3 - (p % 4)));
                  q += 1;
                }
            else if ("UTF16BE" === a || "UTF16LE" === a)
              for (m = 0; m < b.length; m += 1) {
                f = b.charCodeAt(m);
                "UTF16LE" === a && ((n = f & 255), (f = (n << 8) | (f >>> 8)));
                p = q + l;
                for (k = p >>> 2; e.length <= k; ) e.push(0);
                e[k] |= f << (8 * (2 - (p % 4)));
                q += 2;
              }
            return {
              value: e,
              binLen: 8 * q + c
            };
          };
          break;
        case "B64":
          b = L;
          break;
        case "BYTES":
          b = K;
          break;
        case "ARRAYBUFFER":
          try {
            b = new ArrayBuffer(0);
          } catch (g) {
            throw Error("ARRAYBUFFER not supported by this environment");
          }
          b = M;
          break;
        default:
          throw Error("format must be HEX, TEXT, B64, BYTES, or ARRAYBUFFER");
      }
      return b;
    }
  
    function k(c, a) {
      return (c >>> a) | (c << (32 - a));
    }
  
    function N(c, a, b) {
      return (c & a) ^ (~c & b);
    }
  
    function O(c, a, b) {
      return (c & a) ^ (c & b) ^ (a & b);
    }
  
    function P(c) {
      return k(c, 2) ^ k(c, 13) ^ k(c, 22);
    }
  
    function Q(c) {
      return k(c, 6) ^ k(c, 11) ^ k(c, 25);
    }
  
    function R(c) {
      return k(c, 7) ^ k(c, 18) ^ (c >>> 3);
    }
  
    function S(c) {
      return k(c, 17) ^ k(c, 19) ^ (c >>> 10);
    }
  
    function T(c, a) {
      var b = (c & 65535) + (a & 65535);
      return (
        ((((c >>> 16) + (a >>> 16) + (b >>> 16)) & 65535) << 16) | (b & 65535)
      );
    }
  
    function U(c, a, b, g) {
      var d = (c & 65535) + (a & 65535) + (b & 65535) + (g & 65535);
      return (
        ((((c >>> 16) + (a >>> 16) + (b >>> 16) + (g >>> 16) + (d >>> 16)) &
          65535) <<
          16) |
        (d & 65535)
      );
    }
  
    function V(c, a, b, g, d) {
      var f = (c & 65535) + (a & 65535) + (b & 65535) + (g & 65535) + (d & 65535);
      return (
        ((((c >>> 16) +
          (a >>> 16) +
          (b >>> 16) +
          (g >>> 16) +
          (d >>> 16) +
          (f >>> 16)) &
          65535) <<
          16) |
        (f & 65535)
      );
    }
  
    function v(c) {
      var a, b;
      a = [
        3238371032,
        914150663,
        812702999,
        4144912697,
        4290775857,
        1750603025,
        1694076839,
        3204075428
      ];
      b = [
        1779033703,
        3144134277,
        1013904242,
        2773480762,
        1359893119,
        2600822924,
        528734635,
        1541459225
      ];
      switch (c) {
        case "SHA-224":
          c = a;
          break;
        case "SHA-256":
          c = b;
          break;
        case "SHA-384":
          c = [
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l()
          ];
          break;
        case "SHA-512":
          c = [
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l(),
            new l()
          ];
          break;
        default:
          throw Error("Unknown SHA variant");
      }
      return c;
    }
  
    function A(c, a, b) {
      var g,
        d,
        f,
        e,
        h,
        l,
        m,
        n,
        k,
        t,
        p,
        u,
        r,
        v,
        y,
        z,
        A,
        B,
        C,
        D,
        E,
        G,
        w = [],
        F;
      if ("SHA-224" === b || "SHA-256" === b)
        (t = 64),
          (u = 1),
          (G = Number),
          (r = T),
          (v = U),
          (y = V),
          (z = R),
          (A = S),
          (B = P),
          (C = Q),
          (E = O),
          (D = N),
          (F = H);
      else throw Error("Unexpected error in SHA-2 implementation");
      b = a[0];
      g = a[1];
      d = a[2];
      f = a[3];
      e = a[4];
      h = a[5];
      l = a[6];
      m = a[7];
      for (p = 0; p < t; p += 1)
        16 > p
          ? ((k = p * u),
            (n = c.length <= k ? 0 : c[k]),
            (k = c.length <= k + 1 ? 0 : c[k + 1]),
            (w[p] = new G(n, k)))
          : (w[p] = v(A(w[p - 2]), w[p - 7], z(w[p - 15]), w[p - 16])),
          (n = y(m, C(e), D(e, h, l), F[p], w[p])),
          (k = r(B(b), E(b, g, d))),
          (m = l),
          (l = h),
          (h = e),
          (e = r(f, n)),
          (f = d),
          (d = g),
          (g = b),
          (b = r(n, k));
      a[0] = r(b, a[0]);
      a[1] = r(g, a[1]);
      a[2] = r(d, a[2]);
      a[3] = r(f, a[3]);
      a[4] = r(e, a[4]);
      a[5] = r(h, a[5]);
      a[6] = r(l, a[6]);
      a[7] = r(m, a[7]);
      return a;
    }
    var H;
    H = [
      1116352408,
      1899447441,
      3049323471,
      3921009573,
      961987163,
      1508970993,
      2453635748,
      2870763221,
      3624381080,
      310598401,
      607225278,
      1426881987,
      1925078388,
      2162078206,
      2614888103,
      3248222580,
      3835390401,
      4022224774,
      264347078,
      604807628,
      770255983,
      1249150122,
      1555081692,
      1996064986,
      2554220882,
      2821834349,
      2952996808,
      3210313671,
      3336571891,
      3584528711,
      113926993,
      338241895,
      666307205,
      773529912,
      1294757372,
      1396182291,
      1695183700,
      1986661051,
      2177026350,
      2456956037,
      2730485921,
      2820302411,
      3259730800,
      3345764771,
      3516065817,
      3600352804,
      4094571909,
      275423344,
      430227734,
      506948616,
      659060556,
      883997877,
      958139571,
      1322822218,
      1537002063,
      1747873779,
      1955562222,
      2024104815,
      2227730452,
      2361852424,
      2428436474,
      2756734187,
      3204031479,
      3329325298
    ];
    return u;
  })();
  
  export const sha256 = text => {
    const hash = new SHA256("SHA-256", "TEXT");
    hash.update(text);
    return hash.getHash("B64");
  };
  