// Proxy to the stock-portfolio Cloudflare Worker: keeps the read token server-side
export default async () => {
    const token = process.env.PORTFOLIO_READ_TOKEN;
    const baseUrl = process.env.PORTFOLIO_API_URL;

    if (!token || !baseUrl) {
        return Response.json(
            { error: 'Missing configuration: set PORTFOLIO_READ_TOKEN and PORTFOLIO_API_URL' },
            { status: 500 }
        );
    }

    try {
        const res = await fetch(baseUrl.replace(/\/+$/, '') + '/portfolio', {
            headers: { Authorization: `Bearer ${token}` },
        });
        return new Response(await res.text(), {
            status: res.status,
            headers: { 'Content-Type': 'application/json', 'Cache-Control': 'no-store' },
        });
    } catch (err) {
        return Response.json({ error: 'Portfolio service unreachable' }, { status: 502 });
    }
};
