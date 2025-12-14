import { type NextRequest, NextResponse } from "next/server";
import { auth } from "@/config/auth.config";
import { ROUTES } from "./constants/routes";

// 認証が不要なパス（公開パス）
const PUBLIC_PATHS = [
  "/login",
  "/api/auth", // NextAuth.jsの認証エンドポイント
  "/favicon.ico",
  "/_next", // Next.jsの内部リソース
  "/public", // 公開ファイル
] as const;

// 静的ファイルの拡張子
const STATIC_FILE_EXTENSIONS = [
  ".css",
  ".js",
  ".json",
  ".png",
  ".jpg",
  ".jpeg",
  ".gif",
  ".svg",
  ".ico",
  ".woff",
  ".woff2",
  ".ttf",
  ".eot",
] as const;

/**
 * パスが認証不要かどうかを判定する
 */
const isPublicPath = (pathname: string): boolean => {
  // 静的ファイルかどうかをチェック
  const isStaticFile = STATIC_FILE_EXTENSIONS.some((ext) =>
    pathname.endsWith(ext),
  );
  if (isStaticFile) return true;

  // 公開パスに含まれるかどうかをチェック
  return PUBLIC_PATHS.some((path) => {
    if (path.endsWith("*")) {
      // ワイルドカード対応（例: "/api/*"）
      return pathname.startsWith(path.slice(0, -1));
    }
    return pathname.startsWith(path);
  });
};

export default async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  // 公開パスの場合は認証チェックをスキップ
  if (isPublicPath(pathname)) {
    return NextResponse.next();
  }

  try {
    // セッションを取得
    const session = await auth();

    // 認証が必要だがセッションがない場合、ログインページにリダイレクト
    if (!session) {
      const url = new URL(ROUTES.LOGIN, request.url);
      // リダイレクト後に元のページに戻れるよう、callbackUrlを設定
      url.searchParams.set("callbackUrl", pathname);
      return NextResponse.redirect(url);
    }

    // セッションがある場合は通常のレスポンスを返す
    return NextResponse.next();
  } catch (_error) {
    // エラーが発生した場合もログインページにリダイレクト
    const url = new URL(ROUTES.LOGIN, request.url);
    url.searchParams.set("callbackUrl", pathname);
    url.searchParams.set("error", "auth_error");
    return NextResponse.redirect(url);
  }
}

// Proxyを適用するパスを設定
export const config = {
  matcher: [
    /*
     * 以下のパスを除くすべてのパスにProxyを適用:
     * - api/auth (NextAuth.js)
     * - _next/static (静的ファイル)
     * - _next/image (画像最適化ファイル)
     * - favicon.ico (ファビコン)
     */
    "/((?!api/auth|_next/static|_next/image|favicon.ico).*)",
  ],
};
