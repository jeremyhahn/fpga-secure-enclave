module signing_processor (
    input wire clk,
    input wire reset,
    input wire [127:0] message_hash,    // Input message hash
    input wire [1:0] signing_type,      // Signing type (00: RSA, 01: ECDSA, 10: EdDSA)
    input wire full_signature,          // Full signature or partial signature flag
    // Full keys
    input wire [255:0] rsa_full_key,
    input wire [255:0] ecdsa_full_key,
    input wire [255:0] eddsa_full_key,
    // Key shards
    input wire [255:0] rsa_shard,
    input wire [255:0] ecdsa_shard,
    input wire [255:0] eddsa_shard,
    output reg [255:0] signature_out,   // Output signature (full or partial)
    output reg done                     // Signal when the signature is done
);

    wire [255:0] rsa_signature, rsa_partial_sig;
    wire [255:0] ecdsa_signature, ecdsa_partial_sig;
    wire [255:0] eddsa_signature, eddsa_partial_sig;
    wire rsa_done, rsa_partial_done;
    wire ecdsa_done, ecdsa_partial_done;
    wire eddsa_done, eddsa_partial_done;

    // Full signature logic
    rsa_signing_core rsa_signer_full (
        .clk(clk),
        .reset(reset),
        .private_key(rsa_full_key),
        .message_hash(message_hash),
        .signature_out(rsa_signature),
        .done(rsa_done)
    );

    ecdsa_signing_core ecdsa_signer_full (
        .clk(clk),
        .reset(reset),
        .private_key(ecdsa_full_key),
        .message_hash(message_hash),
        .signature_out(ecdsa_signature),
        .done(ecdsa_done)
    );

    eddsa_signing_core eddsa_signer_full (
        .clk(clk),
        .reset(reset),
        .private_key(eddsa_full_key),
        .message_hash(message_hash),
        .signature_out(eddsa_signature),
        .done(eddsa_done)
    );

    // Partial signature logic
    rsa_signing_core rsa_signer_partial (
        .clk(clk),
        .reset(reset),
        .private_key(rsa_shard),
        .message_hash(message_hash),
        .signature_out(rsa_partial_sig),
        .done(rsa_partial_done)
    );

    ecdsa_signing_core ecdsa_signer_partial (
        .clk(clk),
        .reset(reset),
        .private_key(ecdsa_shard),
        .message_hash(message_hash),
        .signature_out(ecdsa_partial_sig),
        .done(ecdsa_partial_done)
    );

    eddsa_signing_core eddsa_signer_partial (
        .clk(clk),
        .reset(reset),
        .private_key(eddsa_shard),
        .message_hash(message_hash),
        .signature_out(eddsa_partial_sig),
        .done(eddsa_partial_done)
    );

    // Choose between full and partial signature based on request
    always @(posedge clk or posedge reset) begin
        if (reset) begin
            signature_out <= 256'b0;
            done <= 0;
        end else if (full_signature) begin
            case (signing_type)
                2'b00: begin
                    signature_out <= rsa_signature;
                    done <= rsa_done;
                end
                2'b01: begin
                    signature_out <= ecdsa_signature;
                    done <= ecdsa_done;
                end
                2'b10: begin
                    signature_out <= eddsa_signature;
                    done <= eddsa_done;
                end
            endcase
        end else begin
            case (signing_type)
                2'b00: begin
                    signature_out <= rsa_partial_sig;
                    done <= rsa_partial_done;
                end
                2'b01: begin
                    signature_out <= ecdsa_partial_sig;
                    done <= ecdsa_partial_done;
                end
                2'b10: begin
                    signature_out <= eddsa_partial_sig;
                    done <= eddsa_partial_done;
                end
            endcase
        end
    end
endmodule