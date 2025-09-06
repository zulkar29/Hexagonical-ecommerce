import { useState } from "react";
import { ChevronDown } from "lucide-react";

const faqs = [
  {
    question: "How do I add a new product?",
    answer: "Go to the Products page and click the 'Add Product' button. Fill in the product details and save.",
  },
  {
    question: "How can I view recent orders?",
    answer: "Navigate to the Orders page to see a list of all recent orders, their status, and details.",
  },
  {
    question: "How do I update my store settings?",
    answer: "Go to Settings from the sidebar. You can update your store name, contact info, and preferences there.",
  },
  {
    question: "How do I contact support?",
    answer: "Open a new support ticket from the Support page or email us at support@shopvendor.com.",
  },
];

export default function Help() {
  const [open, setOpen] = useState(null);
  return (
    <div className="max-w-2xl mx-auto space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-2">Help Center</h1>
        <p className="text-muted-foreground">Find answers to common questions or contact support for more help.</p>
      </div>
      <div className="bg-card rounded-xl shadow border border-border divide-y divide-border">
        {faqs.map((faq, idx) => (
          <div key={idx}>
            <button
              className="w-full flex items-center justify-between px-6 py-4 text-left text-lg font-medium text-foreground focus:outline-none hover:bg-accent/30 transition"
              onClick={() => setOpen(open === idx ? null : idx)}
              aria-expanded={open === idx}
            >
              {faq.question}
              <ChevronDown className={`w-5 h-5 ml-2 transition-transform ${open === idx ? 'rotate-180' : ''}`} />
            </button>
            {open === idx && (
              <div className="px-6 pb-4 text-muted-foreground text-base animate-fade-in">
                {faq.answer}
              </div>
            )}
          </div>
        ))}
      </div>
      <div className="text-center text-sm text-muted-foreground mt-8">
        Still need help? <a href="mailto:support@shopvendor.com" className="text-primary underline">Contact Support</a>
      </div>
    </div>
  );
}