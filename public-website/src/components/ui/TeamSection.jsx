'use client';

import { motion } from 'framer-motion';
import { Linkedin, Twitter, Mail } from 'lucide-react';
import Image from 'next/image';

const TeamMember = ({ member, index }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, delay: index * 0.1 }}
      viewport={{ once: true }}
      className="group"
    >
      <div className="relative overflow-hidden rounded-xl bg-gray-100">
        {member.image ? (
          <Image
            src={member.image}
            alt={member.name}
            width={400}
            height={400}
            className="w-full h-80 object-cover transition-transform duration-300 group-hover:scale-105"
          />
        ) : (
          <div className="w-full h-80 bg-blue-600 flex items-center justify-center">
            <span className="text-white text-4xl font-bold">
              {member.name.split(' ').map(n => n[0]).join('')}
            </span>
          </div>
        )}
        
        {/* Social Links Overlay */}
        <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center">
          <div className="flex space-x-4">
            {member.linkedin && (
              <a
                href={member.linkedin}
                target="_blank"
                rel="noopener noreferrer"
                className="p-2 bg-white/20 rounded-full text-white hover:bg-white/30 transition-colors"
              >
                <Linkedin className="w-5 h-5" />
              </a>
            )}
            {member.twitter && (
              <a
                href={member.twitter}
                target="_blank"
                rel="noopener noreferrer"
                className="p-2 bg-white/20 rounded-full text-white hover:bg-white/30 transition-colors"
              >
                <Twitter className="w-5 h-5" />
              </a>
            )}
            {member.email && (
              <a
                href={`mailto:${member.email}`}
                className="p-2 bg-white/20 rounded-full text-white hover:bg-white/30 transition-colors"
              >
                <Mail className="w-5 h-5" />
              </a>
            )}
          </div>
        </div>
      </div>
      
      <div className="mt-6 text-center">
        <h3 className="text-xl font-semibold text-gray-900 mb-1">{member.name}</h3>
        <p className="text-blue-600 font-medium mb-2">{member.position}</p>
        {member.bio && (
          <p className="text-gray-600 text-sm leading-relaxed">{member.bio}</p>
        )}
        {member.expertise && (
          <div className="mt-3 flex flex-wrap justify-center gap-2">
            {member.expertise.map((skill, skillIndex) => (
              <span
                key={skillIndex}
                className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
              >
                {skill}
              </span>
            ))}
          </div>
        )}
      </div>
    </motion.div>
  );
};

const TeamSection = ({ 
  title = "Meet Our Team",
  subtitle,
  description,
  members = [],
  variant = 'grid',
  showAll = true,
  maxMembers = 8
}) => {
  const displayMembers = showAll ? members : members.slice(0, maxMembers);
  
  const gridCols = {
    grid: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
    large: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
    leadership: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
  };

  return (
    <section className="py-16 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="text-center mb-16">
          {subtitle && (
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6 }}
              viewport={{ once: true }}
              className="text-blue-600 font-semibold text-sm uppercase tracking-wide mb-2"
            >
              {subtitle}
            </motion.p>
          )}
          
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.1 }}
            viewport={{ once: true }}
            className="text-3xl md:text-4xl font-bold text-gray-900 mb-4"
          >
            {title}
          </motion.h2>
          
          {description && (
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              viewport={{ once: true }}
              className="text-lg text-gray-600 max-w-3xl mx-auto"
            >
              {description}
            </motion.p>
          )}
        </div>
        
        {/* Team Grid */}
        {displayMembers.length > 0 ? (
          <div className={`grid ${gridCols[variant]} gap-8`}>
            {displayMembers.map((member, index) => (
              <TeamMember key={member.id || index} member={member} index={index} />
            ))}
          </div>
        ) : (
          <div className={`grid ${gridCols[variant]} gap-8`}>
            {/* Placeholder team members */}
            {Array.from({ length: 6 }).map((_, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: index * 0.1 }}
                viewport={{ once: true }}
                className="group"
              >
                <div className="relative overflow-hidden rounded-xl bg-blue-600">
                  <div className="w-full h-80 flex items-center justify-center">
                    <span className="text-white text-4xl font-bold">
                      {['JD', 'SM', 'AR', 'MJ', 'LK', 'TP'][index]}
                    </span>
                  </div>
                </div>
                
                <div className="mt-6 text-center">
                  <h3 className="text-xl font-semibold text-gray-900 mb-1">
                    {['John Doe', 'Sarah Miller', 'Alex Rodriguez', 'Maria Johnson', 'Luke Kim', 'Taylor Park'][index]}
                  </h3>
                  <p className="text-blue-600 font-medium mb-2">
                    {['CEO & Founder', 'CTO', 'Head of Design', 'VP of Sales', 'Lead Developer', 'Marketing Director'][index]}
                  </p>
                  <p className="text-gray-600 text-sm leading-relaxed">
                    {[
                      'Visionary leader with 15+ years in business strategy and innovation.',
                      'Technology expert specializing in scalable solutions and team leadership.',
                      'Creative director passionate about user experience and design systems.',
                      'Sales strategist focused on building lasting client relationships.',
                      'Full-stack developer with expertise in modern web technologies.',
                      'Marketing professional driving growth through data-driven strategies.'
                    ][index]}
                  </p>
                  <div className="mt-3 flex flex-wrap justify-center gap-2">
                    {[
                      ['Strategy', 'Leadership'],
                      ['React', 'Node.js', 'AWS'],
                      ['UI/UX', 'Figma', 'Design Systems'],
                      ['B2B Sales', 'CRM', 'Negotiation'],
                      ['JavaScript', 'Python', 'DevOps'],
                      ['SEO', 'Analytics', 'Content']
                    ][index].map((skill, skillIndex) => (
                      <span
                        key={skillIndex}
                        className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                      >
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>
              </motion.div>
            ))}
          </div>
        )}
        
        {/* Join Team CTA */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.4 }}
          viewport={{ once: true }}
          className="mt-16 text-center"
        >
          <div className="bg-gray-50 rounded-2xl p-8">
            <h3 className="text-2xl font-bold text-gray-900 mb-4">
              Join Our Growing Team
            </h3>
            <p className="text-gray-600 mb-6 max-w-2xl mx-auto">
              We're always looking for talented individuals who share our passion for innovation and excellence. 
              Explore our open positions and become part of our mission.
            </p>
            <a
              href="/careers"
              className="inline-flex items-center px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors"
            >
              View Open Positions
            </a>
          </div>
        </motion.div>
      </div>
    </section>
  );
};

export default TeamSection;