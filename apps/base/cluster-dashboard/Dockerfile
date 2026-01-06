FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy static files to nginx html directory
COPY src/ /usr/share/nginx/html/

# Use non-root user for security
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Run as non-root user
USER nginx

CMD ["nginx", "-g", "daemon off;"]
