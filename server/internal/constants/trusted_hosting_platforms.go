package constants

// TrustedHostingPlatforms lists well-known legitimate static/app hosting services
// whose subdomains appear in the PSL private section (returning icann=false from
// publicsuffix.PublicSuffix). These are operated by reputable companies and must
// not be penalized as "unregulated" TLDs. Subdomains on these platforms also
// lack per-subdomain NS records (the platform manages the whole zone) and are
// inherently unranked, so rank-0 and missing-NS penalties are suppressed too.
var TrustedHostingPlatforms = map[string]struct{}{
	"github.io":              {}, // GitHub Pages
	"github.dev":             {}, // GitHub dev environments
	"gitlab.io":              {}, // GitLab Pages
	"netlify.app":            {}, // Netlify
	"vercel.app":             {}, // Vercel
	"pages.dev":              {}, // Cloudflare Pages
	"workers.dev":            {}, // Cloudflare Workers
	"web.app":                {}, // Firebase Hosting
	"firebaseapp.com":        {}, // Firebase Hosting (legacy)
	"surge.sh":               {}, // Surge.sh
	"fly.dev":                {}, // Fly.io
	"onrender.com":           {}, // Render
	"herokuapp.com":          {}, // Heroku
	"glitch.me":              {}, // Glitch
	"replit.app":             {}, // Replit
	"bitbucket.io":           {}, // Bitbucket Pages
	"azurewebsites.net":      {}, // Azure App Service
	"azurestaticapps.net":    {}, // Azure Static Web Apps
	"amplifyapp.com":         {}, // AWS Amplify
	"readthedocs.io":         {}, // ReadTheDocs
	"huggingface.co":         {}, // Hugging Face Spaces (has PSL entry)
	"streamlit.app":          {}, // Streamlit Cloud
	"railway.app":            {}, // Railway
	"koyeb.app":              {}, // Koyeb
	"cyclic.app":             {}, // Cyclic
	"deno.dev":               {}, // Deno Deploy
	"val.run":                {}, // Val Town
}
