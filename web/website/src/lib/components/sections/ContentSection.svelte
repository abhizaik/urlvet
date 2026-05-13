<script lang="ts">
  import TooltipIcon from "../TooltipIcon.svelte";
  import type { ContentData } from "../../types";
  export let contentData: ContentData | undefined;
</script>

{#if contentData}
  <section
    id="section-content"
    class="bg-white dark:bg-gray-900/80 border border-gray-300 dark:border-gray-800 rounded-lg p-5 shadow-md hover:shadow-lg hover:scale-[1.01] transition-all scroll-mt-20"
  >
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-base font-semibold text-gray-900 dark:text-white">Page Content Analysis</h3>
      <span
        class="text-[10px] text-gray-500 dark:text-gray-400 uppercase tracking-wide px-2 py-0.5 bg-gray-100 dark:bg-gray-800 rounded"
        >DOM Analysis</span
      >
    </div>

    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-[#424242] dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Page Title:</span>
          <TooltipIcon text="The title of the page as defined in the HTML <title> tag." />
        </div>
        <span class="font-medium text-[#424242] dark:text-white"
          >{contentData.title || "(No Title)"}</span
        >
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Brand Verification:</span>
          <TooltipIcon
            text="Checks if the page content matches well-known brands and verifies if it's hosted on an official domain."
          />
        </div>
        {#if contentData.brand_check?.is_mismatch}
          <span class="text-red-400 font-medium flex items-center gap-1">
            ❌ Brand Mismatch ({contentData.brand_check.brand_found})
          </span>
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1">
            ✅ {contentData.brand_check?.detected_names?.length
              ? "Verified Brands: " + contentData.brand_check.detected_names.join(", ")
              : "No high-value brands detected"}
          </span>
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Forms Detected:</span>
          <TooltipIcon text="Total number of HTML forms found on the page." />
        </div>
        <span class="font-medium text-[#424242] dark:text-white">{contentData.form_count}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Login Form Presence:</span>
          <TooltipIcon
            text="Checks if any forms appear to be for logging in (contain password or username-like fields)."
          />
        </div>
        {#if contentData.has_login_form}
          <span class="text-red-400 font-medium flex items-center gap-1">Detected</span>
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1"
            >None Detected</span
          >
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Payment Form Presence:</span>
          <TooltipIcon
            text="Checks if any forms appear to be for payments (contain credit card, CVV, or billing fields)."
          />
        </div>
        {#if contentData.has_payment_form}
          <span class="text-red-400 font-medium flex items-center gap-1">Detected</span>
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1"
            >None Detected</span
          >
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Personal Info Collection:</span>
          <TooltipIcon
            text="Checks if any forms request sensitive personal info like address, phone, or SSN."
          />
        </div>
        {#if contentData.has_personal_form}
          <span class="text-red-400 font-medium flex items-center gap-1">Detected</span>
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1"
            >None Detected</span
          >
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Hidden Elements:</span>
          <TooltipIcon
            text="Detects forms or iframes that are hidden from view, which can be used for malicious background activities."
          />
        </div>
        {#if contentData.has_hidden_iframe || contentData.forms?.some((f) => f.is_hidden)}
          <span class="text-red-400 font-medium flex items-center gap-1">
            ⚠️ {contentData.has_hidden_iframe ? "Hidden Iframe" : ""}
            {contentData.has_hidden_iframe && contentData.forms?.some((f) => f.is_hidden)
              ? "&"
              : ""}
            {contentData.forms?.some((f) => f.is_hidden) ? "Hidden Form" : ""} Detected
          </span>
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1"
            >None Detected</span
          >
        {/if}
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>Tracking Beacons:</span>
          <TooltipIcon
            text="Detects 1x1 or 0x0 pixel images used for background tracking or verifying email opens."
          />
        </div>
        {#if contentData.has_tracking}
          <span class="text-amber-700 dark:text-amber-400 font-medium flex items-center gap-1"
            >Detected</span
          >
        {:else}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1"
            >None Detected</span
          >
        {/if}
      </div>

      {#if contentData.forms && contentData.forms.length > 0}
        <div class="py-4 last:pb-0">
          <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            Detailed Form Technicals
          </h4>
          <div class="space-y-8">
            {#each contentData.forms as form, i}
              <div
                class="space-y-0 divide-y divide-gray-100 dark:divide-gray-800 border border-gray-300 dark:border-gray-800 rounded-lg bg-gray-50 dark:bg-gray-900/40"
              >
                <div
                  class="bg-gray-100 dark:bg-gray-800/60 px-4 py-2 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center rounded-t-lg"
                >
                  <span class="text-xs font-bold text-blue-400 uppercase">Form #{i + 1}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Submission Method:</span>
                    <TooltipIcon
                      text="The HTTP method used to send data (POST is standard, GET can leak data in URLs)."
                    />
                  </div>
                  <span class="font-mono text-gray-700 dark:text-gray-200 uppercase"
                    >{form.method}</span
                  >
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Submission Endpoint:</span>
                    <TooltipIcon text="The destination URL where the form data will be sent." />
                  </div>
                  <span class="font-mono text-gray-900 dark:text-white break-all"
                    >{form.action || "(Current Page)"}</span
                  >
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Data Flow:</span>
                    <TooltipIcon
                      text="Checks if data is being sent to the same website or an external/unrelated domain."
                    />
                  </div>
                  {#if form.is_external}
                    <span class="text-red-400 font-medium">⚠️ Submits to External Domain</span>
                  {:else}
                    <span class="text-emerald-700 dark:text-emerald-400 font-medium"
                      >✅ Submits to Same Domain</span
                    >
                  {/if}
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Security Analysis:</span>
                    <TooltipIcon
                      text="Automated check for suspicious form properties or sensitive data collection."
                    />
                  </div>
                  <div class="flex flex-wrap gap-2">
                    {#if !form.has_password && !form.has_user_like && !form.has_payment && !form.has_personal && !form.is_hidden}
                      <span class="text-gray-400 italic">No sensitive flags detected</span>
                    {/if}
                    {#if form.is_hidden}
                      <span class="text-red-400 font-bold">👻 HIDDEN FORM</span>
                    {/if}
                    {#if form.has_password}
                      <span class="text-amber-700 dark:text-amber-400 flex items-center gap-1"
                        >🔒 Collects Passwords</span
                      >
                    {/if}
                    {#if form.has_user_like}
                      <span class="text-blue-400 flex items-center gap-1">👤 Identity Fields</span>
                    {/if}
                    {#if form.has_payment}
                      <span class="text-red-400 flex items-center gap-1">💳 Payment Data</span>
                    {/if}
                    {#if form.has_personal}
                      <span class="text-orange-400 flex items-center gap-1">🏠 Personal Info</span>
                    {/if}
                  </div>
                </div>

                {#if form.inputs && form.inputs.length > 0}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 px-4 py-2"
                  >
                    <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                      <span>Detected Data Fields:</span>
                      <TooltipIcon
                        text="Full technical map of input fields found within this form."
                      />
                    </div>
                    <div class="flex flex-col gap-1.5">
                      {#each form.inputs as input}
                        <span
                          class="text-[11px] text-gray-600 dark:text-gray-300 font-mono bg-gray-100 dark:bg-gray-800/50 px-2 py-1 rounded border border-gray-300 dark:border-gray-700/30 break-all"
                        >
                          {input}
                        </span>
                      {/each}
                    </div>
                  </div>
                {/if}

                {#if form.submit_texts && form.submit_texts.length > 0}
                  <div
                    class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                  >
                    <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                      <span>Submission Buttons:</span>
                      <TooltipIcon
                        text="The text labels on buttons that trigger this form's submission."
                      />
                    </div>
                    <div class="flex flex-wrap gap-1 min-w-0">
                      {#each form.submit_texts as text}
                        <span
                          class="px-2 py-0.5 bg-gray-900 text-emerald-400 rounded border border-emerald-900/30 text-xs font-medium break-all"
                        >
                          {text}
                        </span>
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}

      {#if contentData.iframes && contentData.iframes.length > 0}
        <div class="py-4 last:pb-0">
          <h4 class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-4">
            Iframe & Third-Party Elements
          </h4>
          <div class="space-y-8">
            {#each contentData.iframes as iframe, i}
              <div
                class="space-y-0 divide-y divide-gray-100 dark:divide-gray-800 border border-gray-300 dark:border-gray-800 rounded-lg bg-gray-50 dark:bg-gray-900/40"
              >
                <div
                  class="bg-gray-100 dark:bg-gray-800/60 px-4 py-2 border-b border-gray-100 dark:border-gray-800 flex justify-between items-center rounded-t-lg"
                >
                  <span class="text-xs font-bold text-purple-400 uppercase">Iframe #{i + 1}</span>
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Visibility Status:</span>
                    <TooltipIcon
                      text="Indicates if the iframe is visible to the user or hidden in the background."
                    />
                  </div>
                  {#if iframe.is_hidden}
                    <span class="text-red-400 font-bold flex items-center gap-1">👻 Hidden</span>
                  {:else}
                    <span class="text-gray-400 font-medium">Visible</span>
                  {/if}
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Source (URL):</span>
                    <TooltipIcon text="The external URL being loaded into this iframe." />
                  </div>
                  <span class="font-mono text-gray-900 dark:text-white break-all"
                    >{iframe.src || "(No Source)"}</span
                  >
                </div>

                <div
                  class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 px-4 py-2"
                >
                  <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                    <span>Dimensions:</span>
                    <TooltipIcon text="The width and height of the iframe element." />
                  </div>
                  <span class="text-gray-600 dark:text-gray-300 font-mono">
                    {iframe.width || "auto"} x {iframe.height || "auto"}
                  </span>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
