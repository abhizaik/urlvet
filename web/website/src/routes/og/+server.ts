import type { RequestHandler } from '@sveltejs/kit';
import satori from 'satori';
import sharp from 'sharp';

// Cache font in memory after first fetch
let fontData: ArrayBuffer | null = null;

async function getFont(): Promise<ArrayBuffer | null> {
  if (fontData) return fontData;
  try {
    const res = await fetch(
      'https://fonts.gstatic.com/s/inter/v13/UcCO3FwrK3iLTeHuS_fvQtMwCp50KnMw2boKoduKmMEVuFuYAZJhjp-Ek-_EeA.woff'
    );
    fontData = await res.arrayBuffer();
    return fontData;
  } catch {
    return null;
  }
}

const VERDICT_COLORS: Record<string, { bg: string; accent: string; label: string }> = {
  Safe: { bg: '#022c22', accent: '#10b981', label: 'Trusted' },
  Risky: { bg: '#2d0a0a', accent: '#ef4444', label: 'High Risk' },
  Suspicious: { bg: '#1c1408', accent: '#eab308', label: 'Be Cautious' },
};

export const GET: RequestHandler = async ({ url }) => {
  const domain = url.searchParams.get('domain') ?? 'unknown';
  const verdict = url.searchParams.get('v') ?? 'Suspicious';
  const score = parseInt(url.searchParams.get('s') ?? '0', 10);

  const colors = VERDICT_COLORS[verdict] ?? VERDICT_COLORS.Suspicious;
  const font = await getFont();

  const svg = await satori(
    {
      type: 'div',
      props: {
        style: {
          width: '1200px',
          height: '630px',
          background: '#0a0a0f',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          fontFamily: 'Inter, sans-serif',
          padding: '60px',
          gap: '24px',
          position: 'relative',
        },
        children: [
          // Coloured glow blob
          {
            type: 'div',
            props: {
              style: {
                position: 'absolute',
                top: '-80px',
                left: '-80px',
                width: '400px',
                height: '400px',
                borderRadius: '50%',
                background: colors.accent,
                opacity: '0.08',
                filter: 'blur(80px)',
              },
            },
          },
          // Brand
          {
            type: 'div',
            props: {
              style: {
                color: '#6b7280',
                fontSize: '28px',
                letterSpacing: '0.1em',
                textTransform: 'uppercase',
              },
              children: 'SafeSurf',
            },
          },
          // Domain
          {
            type: 'div',
            props: {
              style: {
                color: '#ffffff',
                fontSize: '52px',
                fontWeight: '800',
                letterSpacing: '-0.02em',
                maxWidth: '900px',
                textAlign: 'center',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              },
              children: domain,
            },
          },
          // Verdict row
          {
            type: 'div',
            props: {
              style: { display: 'flex', alignItems: 'center', gap: '16px', marginTop: '8px' },
              children: [
                {
                  type: 'div',
                  props: {
                    style: { color: colors.accent, fontSize: '58px', fontWeight: '800' },
                    children: verdict,
                  },
                },
                {
                  type: 'div',
                  props: {
                    style: {
                      background: colors.bg,
                      border: `1.5px solid ${colors.accent}40`,
                      color: colors.accent,
                      fontSize: '24px',
                      fontWeight: '700',
                      padding: '8px 20px',
                      borderRadius: '999px',
                      letterSpacing: '0.08em',
                      textTransform: 'uppercase',
                    },
                    children: colors.label,
                  },
                },
              ],
            },
          },
          // Score
          {
            type: 'div',
            props: {
              style: { color: '#6b7280', fontSize: '26px', marginTop: '4px' },
              children: `Trust Score: ${score}/100`,
            },
          },
        ],
      },
    },
    {
      width: 1200,
      height: 630,
      fonts: font ? [{ name: 'Inter', data: font, weight: 800, style: 'normal' }] : [],
    }
  );

  const png = await sharp(Buffer.from(svg)).png().toBuffer();

  return new Response(png, {
    headers: {
      'Content-Type': 'image/png',
      'Cache-Control': 'public, max-age=3600',
    },
  });
};
