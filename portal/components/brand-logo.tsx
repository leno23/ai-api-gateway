import Link from "next/link";

type BrandLogoProps = {
  href?: string;
  className?: string;
};

export function BrandLogo({ href = "/", className = "" }: BrandLogoProps) {
  const inner = (
    <span className={`inline-flex items-center gap-2 ${className}`}>
      <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#0f172a] text-sm font-bold text-white">
        蓝
      </span>
      <span className="text-lg font-semibold tracking-tight text-[var(--semi-color-text-0)]">
        蓝移 API
      </span>
    </span>
  );

  if (href) {
    return (
      <Link href={href} className="shrink-0 no-underline">
        {inner}
      </Link>
    );
  }

  return inner;
}
