module rsa_signing_core (
    input wire clk,
    input wire reset,
    input wire [255:0] private_key,    // Full private key for RSA
    input wire [127:0] message_hash,   // Message hash
    input wire partial,                // Use partial signing (threshold)
    input wire [255:0] partial_key,    // Shard of the private key for partial signature
    output reg [255:0] signature_out,  // RSA signature
    output reg done
);

    rsa_full_signer rsa_full (
        .clk(clk),
        .reset(reset),
        .private_key(private_key),
        .message_hash(message_hash),
        .signature(signature_out),
        .done(done)
    );

    rsa_partial_signer rsa_partial (
        .clk(clk),
        .reset(reset),
        .partial_key(partial_key),
        .message_hash(message_hash),
        .signature(signature_out),
        .done(done)
    );

    always @(posedge clk or posedge reset) begin
        if (reset) begin
            signature_out <= 256'b0;
            done <= 0;
        end else if (partial) begin
            // Use partial key signing if indicated
            rsa_partial.partial_key <= partial_key;
            rsa_partial.message_hash <= message_hash;
        end else begin
            // Use full key signing
            rsa_full.private_key <= private_key;
            rsa_full.message_hash <= message_hash;
        end
    end
endmodule