import React, { useState, useEffect } from "react";
import { Home, ShoppingCart, Search } from "lucide-react";

const NotFound = () => {
  const [searchTerm, setSearchTerm] = useState("");
  const [isFloating, setIsFloating] = useState(false);
  
  // Items that have "fallen out of the cart"
  const lostItems = [
    { name: "Stylish Headphones", emoji: "🎧" },
    { name: "Smart Watch", emoji: "⌚" },
    { name: "Wireless Earbuds", emoji: "🎵" },
    { name: "Laptop", emoji: "💻" },
    { name: "Phone Case", emoji: "📱" },
    { name: "Sunglasses", emoji: "😎" }
  ];
  
  // Animation effect
  useEffect(() => {
    const interval = setInterval(() => {
      setIsFloating(prev => !prev);
    }, 1500);
    
    return () => clearInterval(interval);
  }, []);

  // Mock navigation function (substitute for react-router Link)
  const navigateTo = (path) => {
    console.log(`Navigating to: ${path}`);
    // In a real app, this would use your routing system
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-purple-50 to-blue-50 p-4 text-center relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        {lostItems.map((item, index) => (
          <div
            key={index}
            className="absolute text-4xl animate-bounce"
            style={{
              left: `${10 + (index * 15)}%`,
              top: `${(index % 3) * 10}%`,
              animationDelay: `${index * 0.2}s`,
              animationDuration: `${2 + index * 0.5}s`
            }}
          >
            {item.emoji}
          </div>
        ))}
      </div>
      
      {/* Broken shopping cart illustration */}
      <div className={`relative ${isFloating ? 'translate-y-2' : 'translate-y-0'} transition-transform duration-1000 ease-in-out mb-8`}>
        <ShoppingCart className="w-32 h-32 text-blue-400" />
        <div className="absolute top-1/2 left-1/2 w-full h-0.5 bg-red-400 transform -translate-x-1/2 -translate-y-1/2 rotate-45" />
      </div>
      
      <h1 className="text-9xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600 drop-shadow-md">404</h1>
      
      <h2 className="text-3xl font-semibold text-gray-800 mt-4">Oops! Page Not in Stock</h2>
      
      <p className="text-gray-600 mt-4 max-w-md">
        It looks like this item has fallen out of our digital shopping cart!
      </p>
      
      {/* Mock search */}
      <div className="mt-8 relative max-w-md w-full">
        <input
          type="text"
          className="w-full px-4 py-3 rounded-full border-2 border-blue-300 focus:border-blue-500 focus:outline-none pl-12"
          placeholder="What were you looking for?"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <Search className="absolute left-4 top-3.5 w-5 h-5 text-blue-500" />
      </div>
      
      {/* Suggested products */}
      <div className="mt-8 grid grid-cols-3 gap-4 max-w-lg">
        {lostItems.slice(0, 3).map((item, index) => (
          <button 
            key={index} 
            onClick={() => navigateTo("/products")} 
            className="p-4 bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow cursor-pointer"
          >
            <div className="text-3xl mb-2">{item.emoji}</div>
            <div className="text-sm font-medium text-gray-800">{item.name}</div>
          </button>
        ))}
      </div>
      
      <div className="mt-8 flex flex-col sm:flex-row gap-4">
        <button
          onClick={() => navigateTo("/")}
          className="flex items-center justify-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors shadow-md"
        >
          <Home className="w-5 h-5 mr-2" />
          Return Home
        </button>
        
        <button
          onClick={() => navigateTo("/products")}
          className="flex items-center justify-center px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors shadow-md"
        >
          <ShoppingCart className="w-5 h-5 mr-2" />
          Continue Shopping
        </button>
      </div>
      
      <div className="mt-12 text-gray-500 max-w-md">
        <h3 className="font-medium text-lg">Did you know?</h3>
        <p className="mt-2">The 404 error is named after room 404 at CERN, where the World Wide Web was developed. The room reportedly contained all the documentation about the web... until it didn't!</p>
      </div>
    </div>
  );
};

export default NotFound;